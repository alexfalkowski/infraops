package main

import (
	"github.com/alexfalkowski/infraops/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		if err := do.Configure(ctx); err != nil {
			return err
		}

		lt := &do.Project{
			Name:        "lean-thoughts",
			Description: "The lean thoughts domain",
		}

		return do.CreateProject(ctx, lt)
	})
}
