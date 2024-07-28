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
		Version:       "1.210.0",
		Memory:        app.Memory{Min: "128Mi", Max: "256Mi"},
		SecretVolumes: []string{"gh"},
	}

	return app.CreateApp(ctx, a)
}
