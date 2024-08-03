package lean

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "1.220.0",
		Version:       "1.156.0",
		ConfigVersion: "1.10.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}
