package main

import (
	"github.com/alexfalkowski/infraops/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var applications []*app.App

func RegisterApplication(application *app.App) {
	applications = append(applications, application)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, application := range applications {
			if err := app.CreateApp(ctx, application); err != nil {
				return err
			}
		}

		return nil
	})
}
