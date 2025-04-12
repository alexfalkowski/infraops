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
	var args *cv1.ServiceArgs

	if app.IsInternal() {
		args = &cv1.ServiceArgs{
			Metadata: metadata(app, matchLabels(app)),
			Spec: cv1.ServiceSpecArgs{
				Ports: cv1.ServicePortArray{
					servicePort("debug", 6060),
					servicePort("http", 8080),
					servicePort("grpc", 9090),
				},
				Selector: matchLabels(app),
				Type:     pulumi.String("ClusterIP"),
			},
		}
	} else {
		args = &cv1.ServiceArgs{
			Metadata: metadata(app, matchLabels(app)),
			Spec: cv1.ServiceSpecArgs{
				Ports: cv1.ServicePortArray{
					servicePort("http", 8080),
				},
				Selector: matchLabels(app),
				Type:     pulumi.String("ClusterIP"),
			},
		}
	}

	_, err := cv1.NewService(ctx, app.Name, args)

	return err
}

func initContainers(app *App) cv1.ContainerArray {
	if !app.HasConfigVersion() || !app.IsInternal() {
		return nil
	}

	name := initName(app)
	path := configFilePath("konfig", name)
	volumeMounts := cv1.VolumeMountArray{
		cv1.VolumeMountArgs{
			MountPath: path,
			Name:      pulumi.String(name),
			SubPath:   pulumi.String(configFile(name)),
		},
		cv1.VolumeMountArgs{
			Name:      pulumi.String(app.Name),
			MountPath: pulumi.String(configPath(app.Name)),
		},
	}

	for _, s := range app.Secrets {
		volumeMounts = append(volumeMounts, secretVolumeMount(s))
	}

	return cv1.ContainerArray{
		cv1.ContainerArgs{
			Name:            pulumi.String(name),
			Image:           initImage("konfigctl", app.InitVersion),
			ImagePullPolicy: pulumi.String("Always"),
			Args:            pulumi.StringArray{pulumi.String("config")},
			VolumeMounts:    volumeMounts,
			Env: cv1.EnvVarArray{
				cv1.EnvVarArgs{
					Name: pulumi.String("SERVICE_ID"),
					ValueFrom: &cv1.EnvVarSourceArgs{
						FieldRef: &cv1.ObjectFieldSelectorArgs{
							FieldPath: pulumi.String("metadata.uid"),
						},
					},
				},
				cv1.EnvVarArgs{
					Name:  pulumi.String("KONFIG_CONFIG_FILE"),
					Value: path,
				},
				cv1.EnvVarArgs{
					Name:  pulumi.String("KONFIG_APP_CONFIG_FILE"),
					Value: configMatchingFilePath(app.Name),
				},
			},
			Resources: createResources(app),
			SecurityContext: cv1.SecurityContextArgs{
				ReadOnlyRootFilesystem: pulumi.Bool(true),
			},
		},
	}
}

func containers(app *App) cv1.ContainerArray {
	if app.IsInternal() {
		return internalContainer(app)
	}

	return externalContainer(app)
}

func internalContainer(app *App) cv1.ContainerArray {
	volumeMounts := cv1.VolumeMountArray{}

	if app.HasConfigVersion() {
		v := cv1.VolumeMountArgs{
			MountPath: pulumi.String(configPath(app.Name)),
			Name:      pulumi.String(app.Name),
		}
		volumeMounts = append(volumeMounts, v)
	} else {
		v := cv1.VolumeMountArgs{
			MountPath: configMatchingFilePath(app.Name),
			Name:      pulumi.String(app.Name),
			SubPath:   pulumi.String(configFile(app.Name)),
		}
		volumeMounts = append(volumeMounts, v)
	}

	for _, v := range app.Secrets {
		volumeMounts = append(volumeMounts, secretVolumeMount(v))
	}

	envs := cv1.EnvVarArray{}
	envs = append(envs, cv1.EnvVarArgs{
		Name: pulumi.String("SERVICE_ID"),
		ValueFrom: &cv1.EnvVarSourceArgs{
			FieldRef: &cv1.ObjectFieldSelectorArgs{
				FieldPath: pulumi.String("metadata.uid"),
			},
		},
	})
	envs = append(envs, cv1.EnvVarArgs{
		Name:  pulumi.String(strings.ToUpper(app.Name) + "_CONFIG_FILE"),
		Value: configMatchingFilePath(app.Name),
	})

	container := cv1.ContainerArgs{
		Name:            pulumi.String(app.Name),
		Image:           image(app),
		ImagePullPolicy: pulumi.String("Always"),
		Args:            pulumi.StringArray{pulumi.String("server")},
		VolumeMounts:    volumeMounts,
		Env:             addEnvironments(app, envs),
		Ports: cv1.ContainerPortArray{
			cv1.ContainerPortArgs{ContainerPort: pulumi.Int(6060)},
			cv1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
			cv1.ContainerPortArgs{ContainerPort: pulumi.Int(9090)},
		},
		LivenessProbe:  httpProbe("/livez"),
		ReadinessProbe: httpProbe("/readyz"),
		StartupProbe:   tcpProbe(),
		Resources:      createResources(app),
		SecurityContext: cv1.SecurityContextArgs{
			ReadOnlyRootFilesystem: pulumi.Bool(true),
		},
	}

	return cv1.ContainerArray{container}
}

func externalContainer(app *App) cv1.ContainerArray {
	container := cv1.ContainerArgs{
		Name:            pulumi.String(app.Name),
		Image:           image(app),
		ImagePullPolicy: pulumi.String("Always"),
		Env:             addEnvironments(app, cv1.EnvVarArray{}),
		Ports: cv1.ContainerPortArray{
			cv1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
		},
		LivenessProbe:  httpProbe("/"),
		ReadinessProbe: tcpProbe(),
		StartupProbe:   tcpProbe(),
		Resources:      createResources(app),
		SecurityContext: cv1.SecurityContextArgs{
			ReadOnlyRootFilesystem: pulumi.Bool(true),
		},
	}

	return cv1.ContainerArray{container}
}

func addEnvironments(app *App, envs cv1.EnvVarArray) cv1.EnvVarArray {
	for _, env := range app.Environments {
		var arg cv1.EnvVarArgs

		if env.IsSecret() {
			value := strings.TrimPrefix(env.Value, "secret:")
			name, value, _ := strings.Cut(value, "/")

			arg = cv1.EnvVarArgs{
				Name: pulumi.String(env.Name),
				ValueFrom: &cv1.EnvVarSourceArgs{
					SecretKeyRef: &cv1.SecretKeySelectorArgs{
						Name: pulumi.String(name + "-secret"),
						Key:  pulumi.String(value),
					},
				},
			}
		} else {
			arg = cv1.EnvVarArgs{
				Name:  pulumi.String(env.Name),
				Value: pulumi.String(env.Value),
			}
		}

		envs = append(envs, arg)
	}

	return envs
}

func createVolumes(app *App) cv1.VolumeArray {
	volumes := cv1.VolumeArray{}
	if app.IsExternal() {
		return volumes
	}

	if app.HasConfigVersion() {
		n := pulumi.String(initName(app))
		k := cv1.VolumeArgs{
			Name:      n,
			ConfigMap: cv1.ConfigMapVolumeSourceArgs{Name: n},
		}
		s := cv1.VolumeArgs{
			Name:     pulumi.String(app.Name),
			EmptyDir: cv1.EmptyDirVolumeSourceArgs{},
		}
		volumes = append(volumes, k, s)
	} else {
		v := cv1.VolumeArgs{
			Name:      pulumi.String(app.Name),
			ConfigMap: cv1.ConfigMapVolumeSourceArgs{Name: pulumi.String(app.Name)},
		}
		volumes = append(volumes, v)
	}

	for _, v := range app.Secrets {
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

func initImage(name, version string) pulumi.String {
	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:v%s", name, version))
}

func image(app *App) pulumi.String {
	if app.IsInternal() {
		return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:v%s", app.Name, app.Version))
	}

	return pulumi.String(fmt.Sprintf("docker.io/alexfalkowski/%s:%s", app.Name, app.Version))
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

func createResources(app *App) cv1.ResourceRequirementsArgs {
	if !app.HasResources() {
		return cv1.ResourceRequirementsArgs{}
	}

	requests := pulumi.StringMap{}
	limits := pulumi.StringMap{}

	if app.Resources.CPU != nil {
		requests["cpu"] = pulumi.String(app.Resources.CPU.Min)
		limits["cpu"] = pulumi.String(app.Resources.CPU.Max)
	}

	if app.Resources.Memory != nil {
		requests["memory"] = pulumi.String(app.Resources.Memory.Min)
		limits["memory"] = pulumi.String(app.Resources.Memory.Max)
	}

	if app.Resources.Storage != nil {
		requests["ephemeral-storage"] = pulumi.String(app.Resources.Storage.Min)
		limits["ephemeral-storage"] = pulumi.String(app.Resources.Storage.Max)
	}

	return cv1.ResourceRequirementsArgs{
		Requests: requests,
		Limits:   limits,
	}
}

func httpProbe(path string) cv1.ProbeArgs {
	return cv1.ProbeArgs{
		HttpGet: cv1.HTTPGetActionArgs{
			Path: pulumi.String(path),
			Port: pulumi.Int(8080),
			HttpHeaders: cv1.HTTPHeaderArray{
				cv1.HTTPHeaderArgs{
					Name:  pulumi.String("Content-Type"),
					Value: pulumi.String("application/json"),
				},
			},
		},
		InitialDelaySeconds: pulumi.Int(5),
		PeriodSeconds:       pulumi.Int(10),
		TimeoutSeconds:      pulumi.Int(30),
	}
}

func tcpProbe() cv1.ProbeArgs {
	return cv1.ProbeArgs{
		TcpSocket: cv1.TCPSocketActionArgs{
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds: pulumi.Int(5),
		PeriodSeconds:       pulumi.Int(10),
		TimeoutSeconds:      pulumi.Int(30),
	}
}
