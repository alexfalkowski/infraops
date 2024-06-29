package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "28c679dc-5924-47e8-ac48-73cd842ba5cd",
		Name:          "standort",
		Domain:        "lean-thoughts.com",
		InitVersion:   "1.177.0",
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
