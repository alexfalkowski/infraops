package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createKonfig(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{Name: "konfig", Version: app.KonfigVersion, SecretVolumes: []string{"gh"}})
}

func createStandort(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{Name: "standort", Version: "2.92.1", ConfigVersion: "1.7.0"})
}

func createBezeichner(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{Name: "bezeichner", Version: "1.93.3", ConfigVersion: "1.6.0"})
}
