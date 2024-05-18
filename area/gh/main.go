package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

var fns = []createFn{
	createSite, createAppConfig,
	createInfraOps, createDocker, createBin,
	createNonnative, createGoHealth, createGoService,
	createGoServiceTemplate, createGoClientTemplate,
	createStatus, createStandort, createAuth,
	createKonfig, createMigrieren, createBezeichner,
	createServiceControl,
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
