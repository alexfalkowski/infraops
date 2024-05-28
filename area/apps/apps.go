package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "1115c470-ccc9-4daf-8459-ef1e19c40afe",
		Name:          "konfig",
		Version:       "1.133.0",
		Memory:        app.Memory{Min: "128Mi", Max: "256Mi"},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "28c679dc-5924-47e8-ac48-73cd842ba5cd",
		Name:          "standort",
		InitVersion:   "1.133.0",
		Version:       "2.95.0",
		ConfigVersion: "1.7.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		InitVersion:   "1.133.0",
		Version:       "1.96.0",
		ConfigVersion: "1.6.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}
