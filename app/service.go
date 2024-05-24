package app

import (
	"strings"

	av1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//nolint:funlen
func createDeployment(ctx *pulumi.Context, app *App) error {
	volumes := cv1.VolumeArray{
		cv1.VolumeArgs{
			Name:      pulumi.String(app.Name),
			ConfigMap: cv1.ConfigMapVolumeSourceArgs{Name: pulumi.String(app.Name)},
		},
		secretVolume("otlp"),
	}

	volumeMounts := cv1.VolumeMountArray{
		cv1.VolumeMountArgs{
			MountPath: configPath(app),
			Name:      pulumi.String(app.Name),
			SubPath:   pulumi.String(app.Config),
		},
		secretVolumeMount("otlp"),
	}

	for _, v := range app.SecretVolumes {
		volumes = append(volumes, secretVolume(v))
		volumeMounts = append(volumeMounts, secretVolumeMount(v))
	}

	args := &av1.DeploymentArgs{
		Metadata: metadata(app),
		Spec: av1.DeploymentSpecArgs{
			Selector: mv1.LabelSelectorArgs{MatchLabels: labels(app)},
			Replicas: pulumi.Int(3),
			Strategy: av1.DeploymentStrategyArgs{
				RollingUpdate: av1.RollingUpdateDeploymentArgs{
					MaxSurge:       pulumi.Int(1),
					MaxUnavailable: pulumi.Int(1),
				},
			},
			Template: cv1.PodTemplateSpecArgs{
				Metadata: mv1.ObjectMetaArgs{Labels: labels(app)},
				Spec: cv1.PodSpecArgs{
					ServiceAccountName: pulumi.String(app.Name),
					SecurityContext:    podSecurity(),
					Volumes:            volumes,
					Containers: cv1.ContainerArray{
						cv1.ContainerArgs{
							Name:            pulumi.String(app.Name),
							Image:           image(app),
							ImagePullPolicy: pulumi.String("Always"),
							Args:            pulumi.StringArray{pulumi.String("server")},
							VolumeMounts:    volumeMounts,
							Env: cv1.EnvVarArray{
								cv1.EnvVarArgs{
									Name:  pulumi.String(strings.ToUpper(app.Name)) + "_CONFIG_FILE",
									Value: configPath(app),
								},
							},
							Ports: cv1.ContainerPortArray{
								cv1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
								cv1.ContainerPortArgs{ContainerPort: pulumi.Int(9090)},
							},
							LivenessProbe:  probe("/livez"),
							ReadinessProbe: probe("/readyz"),
							Resources: cv1.ResourceRequirementsArgs{
								Requests: resourceRequirement("125m", "1Gi", "64Mi"),
								Limits:   resourceRequirement("250m", "2Gi", "128Mi"),
							},
							SecurityContext: cv1.SecurityContextArgs{
								ReadOnlyRootFilesystem: pulumi.Bool(true),
							},
						},
					},
				},
			},
		},
	}

	_, err := av1.NewDeployment(ctx, app.Name, args)

	return err
}

func createService(ctx *pulumi.Context, app *App) error {
	args := &cv1.ServiceArgs{
		Metadata: mv1.ObjectMetaArgs{
			Name:      pulumi.String(app.Name),
			Namespace: pulumi.String(app.Name),
			Labels:    labels(app),
		},
		Spec: cv1.ServiceSpecArgs{
			Ports:    cv1.ServicePortArray{servicePort("http", 8080), servicePort("grpc", 9090)},
			Selector: labels(app),
			Type:     pulumi.String("ClusterIP"),
		},
	}
	_, err := cv1.NewService(ctx, app.Name, args)

	return err
}
