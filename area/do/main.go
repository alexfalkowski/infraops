package main

import (
	"github.com/alexfalkowski/infraops/internal/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		if err := do.Configure(ctx); err != nil {
			return err
		}

		config, err := do.ReadConfiguration("do.pbtxt")
		if err != nil {
			return err
		}

		for _, project := range config.GetProjects() {
			if err := do.CreateProject(ctx, do.ConvertProject(project)); err != nil {
				return err
			}
		}

		return nil
	})
}
