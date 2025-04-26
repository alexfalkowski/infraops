package app

import (
	"github.com/alexfalkowski/infraops/internal/inputs"
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
					MaxSurge:       inputs.One,
					MaxUnavailable: inputs.One,
				},
			},
			Template: cv1.PodTemplateSpecArgs{
				Metadata: mv1.ObjectMetaArgs{
					Labels: merge(matchLabels(app), deploymentLabels(app)),
				},
				Spec: cv1.PodSpecArgs{
					ServiceAccountName: pulumi.String(app.Name),
					SecurityContext:    podSecurity(),
					Volumes:            createSecretVolumes(app),
					InitContainers:     initContainers(app),
					Containers:         containers(app),
				},
			},
		},
	}

	_, err := av1.NewDeployment(ctx, app.Name, args)

	return err
}
