package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{
		Name: "konfig", Version: app.KonfigVersion,
		SecretVolumes: []string{"gh"},
	})
}

func createStandort(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{
		Name: "standort", Version: "v2.92.1", ConfigVersion: "v1.7.0",
	})
}
