package main

import (
	"github.com/alexfalkowski/infraops/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := app.Read("apps.txtpb")
		if err != nil {
			return err
		}

		for _, application := range config.GetApplications() {
			if err := app.Create(ctx, app.Convert(application)); err != nil {
				return err
			}
		}

		return nil
	})
}
