package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		fns := []createFn{
			createDocker,
			createAppConfig,
			createTemplate,
			createBin,
			createNonnative,
			createGoHealth,
			createGoService,
		}

		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
