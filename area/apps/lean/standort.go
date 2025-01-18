package lean

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "28c679dc-5924-47e8-ac48-73cd842ba5cd",
		Name:          "standort",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "0.107.0",
		Version:       "2.254.0",
		ConfigVersion: "1.12.1",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp"},
	}

	return app.CreateApp(ctx, a)
}
