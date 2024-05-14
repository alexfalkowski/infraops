package main

import (
	"github.com/alexfalkowski/infraops/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		lt := &do.Project{
			Name:        "lean-thoughts",
			Description: "The lean thoughts domain",
		}

		err := do.CreateProject(ctx, lt)
		if err != nil {
			return err
		}

		t := &do.Project{
			Name:        "test-project",
			Description: "The test project",
		}

		err = do.CreateProject(ctx, t)

		return err
	})
}
