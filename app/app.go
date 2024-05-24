package app

import (
	"fmt"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// App to be created.
type App struct {
	Name          string
	Config        string
	Version       string
	SecretVolumes []string
}

// CreateApp in k8s.
func CreateApp(ctx *pulumi.Context, app *App) error {
	err := createServiceAccount(ctx, app)
	if err != nil {
		return err
	}

	err = createNetworkPolicy(ctx, app)
	if err != nil {
		return err
	}

	err = createConfigMap(ctx, app)
	if err != nil {
		return err
	}

	err = createDeployment(ctx, app)
	if err != nil {
		return err
	}

	err = createService(ctx, app)
	if err != nil {
		return err
	}

	return nil
}

func metadata(app *App) mv1.ObjectMetaArgs {
	return mv1.ObjectMetaArgs{
		Name:      pulumi.String(app.Name),
		Namespace: pulumi.String(app.Name),
	}
}

func image(app *App) pulumi.String {
	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:%s", app.Name, app.Version))
}

func labels(app *App) pulumi.StringMap {
	return pulumi.StringMap{"app": pulumi.String(app.Name)}
}

func secretVolume(name string) cv1.VolumeArgs {
	return cv1.VolumeArgs{
		Name:   pulumi.String(name),
		Secret: cv1.SecretVolumeSourceArgs{SecretName: pulumi.String(name + "-secret")},
	}
}

func secretVolumeMount(name string) cv1.VolumeMountArgs {
	return cv1.VolumeMountArgs{
		Name:      pulumi.String(name),
		MountPath: pulumi.String("/etc/secrets/" + name),
		ReadOnly:  pulumi.Bool(true),
	}
}

func configPath(app *App) pulumi.String {
	return pulumi.String(fmt.Sprintf("/etc/%s/%s", app.Name, app.Config))
}

func podSecurity() cv1.PodSecurityContextArgs {
	return cv1.PodSecurityContextArgs{
		FsGroup:      pulumi.Int(2000),
		RunAsNonRoot: pulumi.Bool(true),
		RunAsUser:    pulumi.Int(10001),
		RunAsGroup:   pulumi.Int(10001),
	}
}

func servicePort(name string, port int) cv1.ServicePortArgs {
	return cv1.ServicePortArgs{
		AppProtocol: pulumi.String("TCP"),
		Name:        pulumi.String(name),
		Port:        pulumi.Int(port),
		TargetPort:  pulumi.Int(port),
	}
}

func resourceRequirement(cpu, storage, memory string) pulumi.StringMap {
	return pulumi.StringMap{
		"cpu":               pulumi.String(cpu),
		"ephemeral-storage": pulumi.String(storage),
		"memory":            pulumi.String(memory),
	}
}

func probe(path string) cv1.ProbeArgs {
	return cv1.ProbeArgs{
		HttpGet: cv1.HTTPGetActionArgs{
			Path: pulumi.String(path),
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds: pulumi.Int(3),
		PeriodSeconds:       pulumi.Int(3),
	}
}
