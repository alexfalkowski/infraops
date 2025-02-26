package lean

import (
	"github.com/alexfalkowski/infraops/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "0.192.0",
		Version:       "1.348.0",
		ConfigVersion: "1.14.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp"},
	}

	return app.CreateApp(ctx, a)
}
