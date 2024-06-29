package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		Domain:        "lean-thoughts.com",
		InitVersion:   "1.177.0",
		Version:       "1.121.0",
		ConfigVersion: "1.8.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	if err := app.CreateApp(ctx, a); err != nil {
		return err
	}

	d, err := a.Probe(ctx, "v1/generate", `{ "application": "uuid", "count": 10 }`)
	if err != nil {
		return err
	}

	return ctx.Log.Info(d, nil)
}
