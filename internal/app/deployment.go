package app

import (
	"fmt"

	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	av1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDeployment(ctx *pulumi.Context, account *cv1.ServiceAccount, config *cv1.ConfigMap, app *App, claims []pulumi.Resource) error {
	var strategy av1.DeploymentStrategyArgs
	switch app.Strategy {
	case "", "rolling":
		strategy.RollingUpdate = av1.RollingUpdateDeploymentArgs{
			MaxSurge:       inputs.One,
			MaxUnavailable: inputs.One,
		}
	case "recreate":
		strategy.Type = pulumi.String("Recreate")
	default:
		return fmt.Errorf("application %s: unsupported strategy %q", app.Name, app.Strategy)
	}

	volumes := createSecretVolumes(config.Metadata.Name().Elem(), app)
	for _, storage := range app.Storage {
		name := pulumi.String(storage.claimName(app.Name))
		volumes = append(volumes, cv1.VolumeArgs{
			Name: name,
			PersistentVolumeClaim: cv1.PersistentVolumeClaimVolumeSourceArgs{
				ClaimName: name,
				ReadOnly:  inputs.No,
			},
		})
	}

	m := metadata(app, matchLabels(app), deploymentLabels(app))
	m.Annotations = deploymentAnnotations(app)
	args := &av1.DeploymentArgs{
		Metadata: m,
		Spec: av1.DeploymentSpecArgs{
			Selector: mv1.LabelSelectorArgs{MatchLabels: matchLabels(app)},
			Replicas: pulumi.Int(int(app.Replicas)),
			Strategy: strategy,
			Template: cv1.PodTemplateSpecArgs{
				Metadata: mv1.ObjectMetaArgs{
					Labels: merge(matchLabels(app), recommendedLabels(app), deploymentLabels(app)),
				},
				Spec: cv1.PodSpecArgs{
					ServiceAccountName:            account.Metadata.Name().Elem(),
					SecurityContext:               podSecurity(),
					TerminationGracePeriodSeconds: pulumi.Int(35),
					Volumes:                       volumes,
					Containers:                    containers(app),
				},
			},
		},
	}
	_, err := av1.NewDeployment(ctx, app.Name, args, pulumi.DependsOn(claims))

	return err
}
