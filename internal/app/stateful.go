package app

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	av1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createStatefulSet(ctx *pulumi.Context, app *App) error {
	m := metadata(app, deploymentLabels(app))
	m.Annotations = deploymentAnnotations(app)

	args := &av1.StatefulSetArgs{
		Metadata: m,
		Spec: av1.StatefulSetSpecArgs{
			Selector:    mv1.LabelSelectorArgs{MatchLabels: matchLabels(app)},
			Replicas:    pulumi.Int(3),
			ServiceName: pulumi.String(app.Name),
			UpdateStrategy: av1.StatefulSetUpdateStrategyArgs{
				RollingUpdate: av1.RollingUpdateStatefulSetStrategyArgs{
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
					Containers:         containers(app),
				},
			},
		},
	}

	_, err := av1.NewStatefulSet(ctx, app.Name, args)

	return err
}
