package main

import (
	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

var fns = []createFn{createKonfig}

func createKonfig(ctx *pulumi.Context) error {
	return app.CreateApp(ctx, &app.App{
		Name: "konfig", Config: "server.yaml", Version: "v1.131.1",
		SecretVolumes: []string{"gh"},
	})
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
