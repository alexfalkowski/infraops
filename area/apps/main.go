package main

import (
	"github.com/alexfalkowski/infraops/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := app.ReadConfiguration("apps.pbtxt")
		if err != nil {
			return err
		}

		for _, application := range config.GetApplications() {
			if err := app.CreateApplication(ctx, app.ConvertApplication(application)); err != nil {
				return err
			}
		}

		return nil
	})
}
