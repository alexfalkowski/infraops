package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		Name:          "konfig",
		Version:       "1.132.1",
		Memory:        app.Memory{Min: "128Mi", Max: "256Mi"},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		Name:          "standort",
		InitVersion:   "1.132.1",
		Version:       "2.94.1",
		ConfigVersion: "1.7.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		Name:          "bezeichner",
		InitVersion:   "1.132.1",
		Version:       "1.95.1",
		ConfigVersion: "1.6.0",
		Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
	}

	return app.CreateApp(ctx, a)
}
