package app

import (
	"fmt"
	"strings"

	av1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDeployment(ctx *pulumi.Context, app *App) error {
	m := metadata(app, deploymentLabels(app))
	m.Annotations = deploymentAnnotations(app)

	args := &av1.DeploymentArgs{
		Metadata: m,
		Spec: av1.DeploymentSpecArgs{
			Selector: mv1.LabelSelectorArgs{MatchLabels: matchLabels(app)},
			Replicas: pulumi.Int(3),
			Strategy: av1.DeploymentStrategyArgs{
				RollingUpdate: av1.RollingUpdateDeploymentArgs{
					MaxSurge:       pulumi.Int(1),
					MaxUnavailable: pulumi.Int(1),
				},
			},
			Template: cv1.PodTemplateSpecArgs{
				Metadata: mv1.ObjectMetaArgs{
					Labels: merge(matchLabels(app), deploymentLabels(app)),
				},
				Spec: cv1.PodSpecArgs{
					ServiceAccountName: pulumi.String(app.Name),
					SecurityContext:    podSecurity(),
					Volumes:            createVolumes(app),
					InitContainers:     initContainers(app),
					Containers:         containers(app),
				},
			},
		},
	}

	_, err := av1.NewDeployment(ctx, app.Name, args)

	return err
}

func createService(ctx *pulumi.Context, app *App) error {
	args := &cv1.ServiceArgs{
		Metadata: metadata(app, matchLabels(app)),
		Spec: cv1.ServiceSpecArgs{
			Ports:    cv1.ServicePortArray{servicePort("http", 8080), servicePort("grpc", 9090)},
			Selector: matchLabels(app),
			Type:     pulumi.String("ClusterIP"),
		},
	}
	_, err := cv1.NewService(ctx, app.Name, args)

	return err
}

func initContainers(app *App) cv1.ContainerArray {
	if app.ConfigVersion == "" {
		return nil
	}

	path := configFullFilePath("konfig")

	volumeMounts := cv1.VolumeMountArray{
		cv1.VolumeMountArgs{
			MountPath: path,
			Name:      pulumi.String("konfig"),
			SubPath:   pulumi.String(configFile("konfig")),
		},
		cv1.VolumeMountArgs{
			Name:      pulumi.String(app.Name),
			MountPath: pulumi.String(configPath(app.Name)),
		},
		secretVolumeMount("konfig"),
		secretVolumeMount("otlp"),
	}

	return cv1.ContainerArray{
		cv1.ContainerArgs{
			Name:            pulumi.String(app.Name + "-init"),
			Image:           image("konfig", app.InitVersion),
			ImagePullPolicy: pulumi.String("Always"),
			Args:            pulumi.StringArray{pulumi.String("config")},
			VolumeMounts:    volumeMounts,
			Env: cv1.EnvVarArray{
				cv1.EnvVarArgs{
					Name:  pulumi.String("KONFIG_CONFIG_FILE"),
					Value: path,
				},
				cv1.EnvVarArgs{
					Name:  pulumi.String("KONFIG_APP_CONFIG_FILE"),
					Value: configFullFilePath(app.Name),
				},
			},
			Resources: cv1.ResourceRequirementsArgs{
				Requests: resourceRequirement("125m", "1Gi", "64Mi"),
				Limits:   resourceRequirement("250m", "2Gi", "128Mi"),
			},
			SecurityContext: cv1.SecurityContextArgs{
				ReadOnlyRootFilesystem: pulumi.Bool(true),
			},
		},
	}
}

func containers(app *App) cv1.ContainerArray {
	volumeMounts := cv1.VolumeMountArray{}

	if app.ConfigVersion != "" {
		v := cv1.VolumeMountArgs{
			MountPath: pulumi.String(configPath(app.Name)),
			Name:      pulumi.String(app.Name),
		}
		volumeMounts = append(volumeMounts, v)
	} else {
		v := cv1.VolumeMountArgs{
			MountPath: configFullFilePath(app.Name),
			Name:      pulumi.String(app.Name),
			SubPath:   pulumi.String(configFile(app.Name)),
		}
		volumeMounts = append(volumeMounts, v)
	}

	volumeMounts = append(volumeMounts, secretVolumeMount("otlp"))

	for _, v := range app.SecretVolumes {
		volumeMounts = append(volumeMounts, secretVolumeMount(v))
	}

	return cv1.ContainerArray{
		cv1.ContainerArgs{
			Name:            pulumi.String(app.Name),
			Image:           image(app.Name, app.Version),
			ImagePullPolicy: pulumi.String("Always"),
			Args:            pulumi.StringArray{pulumi.String("server")},
			VolumeMounts:    volumeMounts,
			Env: cv1.EnvVarArray{
				cv1.EnvVarArgs{
					Name:  pulumi.String(strings.ToUpper(app.Name) + "_CONFIG_FILE"),
					Value: configFullFilePath(app.Name),
				},
			},
			Ports: cv1.ContainerPortArray{
				cv1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
				cv1.ContainerPortArgs{ContainerPort: pulumi.Int(9090)},
			},
			LivenessProbe:  httpProbe("/livez"),
			ReadinessProbe: httpProbe("/readyz"),
			StartupProbe:   tcpProbe(),
			Resources: cv1.ResourceRequirementsArgs{
				Requests: resourceRequirement("125m", "1Gi", app.Memory.Min),
				Limits:   resourceRequirement("250m", "2Gi", app.Memory.Max),
			},
			SecurityContext: cv1.SecurityContextArgs{
				ReadOnlyRootFilesystem: pulumi.Bool(true),
			},
		},
	}
}

func createVolumes(app *App) cv1.VolumeArray {
	volumes := cv1.VolumeArray{}

	if app.ConfigVersion != "" {
		k := cv1.VolumeArgs{
			Name:      pulumi.String("konfig"),
			ConfigMap: cv1.ConfigMapVolumeSourceArgs{Name: pulumi.String("konfig")},
		}
		s := cv1.VolumeArgs{
			Name:     pulumi.String(app.Name),
			EmptyDir: cv1.EmptyDirVolumeSourceArgs{},
		}
		volumes = append(volumes, k, s, secretVolume("konfig"))
	} else {
		v := cv1.VolumeArgs{
			Name:      pulumi.String(app.Name),
			ConfigMap: cv1.ConfigMapVolumeSourceArgs{Name: pulumi.String(app.Name)},
		}
		volumes = append(volumes, v)
	}

	volumes = append(volumes, secretVolume("otlp"))

	for _, v := range app.SecretVolumes {
		volumes = append(volumes, secretVolume(v))
	}

	return volumes
}

func secretVolume(name string) cv1.VolumeArgs {
	n := pulumi.String(name + "-secret")

	return cv1.VolumeArgs{
		Name:   n,
		Secret: cv1.SecretVolumeSourceArgs{SecretName: n},
	}
}

func secretVolumeMount(name string) cv1.VolumeMountArgs {
	return cv1.VolumeMountArgs{
		Name:      pulumi.String(name + "-secret"),
		MountPath: pulumi.String("/etc/secrets/" + name),
		ReadOnly:  pulumi.Bool(true),
	}
}

func image(name, version string) pulumi.String {
	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:v%s", name, version))
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

func httpProbe(path string) cv1.ProbeArgs {
	return cv1.ProbeArgs{
		HttpGet: cv1.HTTPGetActionArgs{
			Path: pulumi.String(path),
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds: pulumi.Int(5),
		PeriodSeconds:       pulumi.Int(5),
	}
}

func tcpProbe() cv1.ProbeArgs {
	return cv1.ProbeArgs{
		TcpSocket: cv1.TCPSocketActionArgs{
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds: pulumi.Int(5),
		PeriodSeconds:       pulumi.Int(5),
	}
}
