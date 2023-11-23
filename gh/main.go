package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		if err := createDocker(ctx); err != nil {
			return err
		}

		if err := createAppConfig(ctx); err != nil {
			return err
		}

		if err := createTemplate(ctx); err != nil {
			return err
		}

		return nil
	})
}
