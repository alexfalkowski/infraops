package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		Name:          "konfig",
		Version:       "1.133.0",
		Memory:        app.Memory{Min: "128Mi", Max: "256Mi"},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
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
		Name:          "bezeichner",
		InitVersion:   "1.133.0",
		Version:       "1.96.0",
		ConfigVersion: "1.6.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}
