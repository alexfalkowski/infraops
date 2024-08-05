package lean

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		ID:        "1115c470-ccc9-4daf-8459-ef1e19c40afe",
		Name:      "konfig",
		Namespace: "lean",
		Domain:    "lean-thoughts.com",
		Version:   "1.224.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "250m", Max: "500m"},
			Memory:  &app.Range{Min: "128Mi", Max: "256Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}
