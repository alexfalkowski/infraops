package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	a := &app.App{
		Name: "konfig", Version: app.KonfigVersion,
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}

func createStandort(ctx *pulumi.Context) error {
	a := &app.App{Name: "standort", Version: "2.94.0", ConfigVersion: "1.7.0"}

	return app.CreateApp(ctx, a)
}

func createBezeichner(ctx *pulumi.Context) error {
	a := &app.App{Name: "bezeichner", Version: "1.95.0", ConfigVersion: "1.6.0"}

	return app.CreateApp(ctx, a)
}
