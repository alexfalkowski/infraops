package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const secretSuffix = "-secret"

func createSecretVolumes(app *App) cv1.VolumeArray {
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
	n := pulumi.String(name + secretSuffix)

	return cv1.VolumeArgs{
		Name:   n,
		Secret: cv1.SecretVolumeSourceArgs{SecretName: n},
	}
}

func secretVolumeMount(name string) cv1.VolumeMountArgs {
	return cv1.VolumeMountArgs{
		Name:      pulumi.String(name + secretSuffix),
		MountPath: pulumi.String("/etc/secrets/" + name),
		ReadOnly:  pulumi.Bool(true),
	}
}
