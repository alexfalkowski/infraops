package app

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	pv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/policy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createPodDisruptionBudget(ctx *pulumi.Context, app *App) error {
	args := &pv1.PodDisruptionBudgetArgs{
		Metadata: metadata(app),
		Spec: pv1.PodDisruptionBudgetSpecArgs{
			MaxUnavailable: inputs.One,
			Selector: v1.LabelSelectorArgs{
				MatchLabels: matchLabels(app),
			},
		},
	}

	_, err := pv1.NewPodDisruptionBudget(ctx, app.Name, args)

	return err
}
