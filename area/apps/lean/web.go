package lean

import (
	"github.com/alexfalkowski/infraops/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createWeb(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "b46608ae-950a-46bb-b37a-4dfe68a95b52",
		Name:          "web",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "0.202.0",
		Version:       "0.229.0",
		ConfigVersion: "1.6.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp"},
	}

	return app.CreateApp(ctx, a)
}
