package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "1115c470-ccc9-4daf-8459-ef1e19c40afe",
		Name:          "konfig",
		Domain:        "lean-thoughts.com",
		Version:       "1.174.0",
		Memory:        app.Memory{Min: "128Mi", Max: "256Mi"},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "28c679dc-5924-47e8-ac48-73cd842ba5cd",
		Name:          "standort",
		Domain:        "lean-thoughts.com",
		InitVersion:   "1.174.0",
		Version:       "2.117.0",
		ConfigVersion: "1.9.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	if err := app.CreateApp(ctx, a); err != nil {
		return err
	}

	d, err := a.Probe(ctx, "v2/location", "{}")
	if err != nil {
		return err
	}

	return ctx.Log.Info(d, nil)
}

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		Domain:        "lean-thoughts.com",
		InitVersion:   "1.174.0",
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
