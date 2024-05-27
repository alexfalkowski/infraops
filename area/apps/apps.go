package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		Name:          "konfig",
		Version:       "1.132.1",
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{
		Name:    "standort",
		Version: "2.94.1", InitVersion: "1.132.1",
		ConfigVersion: "1.7.0",
	}

	return app.CreateApp(ctx, a)
}

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{
		Name:    "bezeichner",
		Version: "1.95.1", InitVersion: "1.132.1",
		ConfigVersion: "1.6.0",
	}

	return app.CreateApp(ctx, a)
}
