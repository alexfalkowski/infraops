package app

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func containers(app *App) cv1.ContainerArray {
	if app.IsInternal() {
		return internalContainer(app)
	}

	return externalContainer(app)
}

func internalContainer(app *App) cv1.ContainerArray {
	volumeMounts := cv1.VolumeMountArray{
		cv1.VolumeMountArgs{
			MountPath: pulumi.String(configMatchingFilePath(app.Name)),
			Name:      pulumi.String(app.Name),
			SubPath:   pulumi.String(configFile(app.Name)),
		},
	}

	for _, v := range app.Secrets {
		volumeMounts = append(volumeMounts, secretVolumeMount(v))
	}

	envs := cv1.EnvVarArray{serviceID()}
	container := cv1.ContainerArgs{
		Name:            pulumi.String(app.Name),
		Image:           image(app),
		ImagePullPolicy: inputs.Always,
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
			ReadOnlyRootFilesystem:   inputs.Yes,
			AllowPrivilegeEscalation: inputs.No,
		},
	}

	return cv1.ContainerArray{container}
}

func externalContainer(app *App) cv1.ContainerArray {
	container := cv1.ContainerArgs{
		Name:            pulumi.String(app.Name),
		Image:           image(app),
		ImagePullPolicy: inputs.Always,
		Env:             addEnvironments(app, cv1.EnvVarArray{}),
		Ports: cv1.ContainerPortArray{
			cv1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
		},
		LivenessProbe:  httpProbe("/"),
		ReadinessProbe: tcpProbe(),
		StartupProbe:   tcpProbe(),
		Resources:      createResources(app),
		SecurityContext: cv1.SecurityContextArgs{
			ReadOnlyRootFilesystem:   inputs.Yes,
			AllowPrivilegeEscalation: inputs.No,
		},
	}

	return cv1.ContainerArray{container}
}

func serviceID() cv1.EnvVarArgs {
	return cv1.EnvVarArgs{
		Name: pulumi.String("SERVICE_ID"),
		ValueFrom: &cv1.EnvVarSourceArgs{
			FieldRef: &cv1.ObjectFieldSelectorArgs{
				FieldPath: pulumi.String("metadata.uid"),
			},
		},
	}
}
