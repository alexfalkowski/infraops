package main

import (
	"github.com/alexfalkowski/infraops/internal/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := gh.ReadConfiguration("gh.pbtxt")
		if err != nil {
			return err
		}

		for _, repository := range config.GetRepositories() {
			if err := gh.CreateRepository(ctx, gh.ConvertRepository(repository)); err != nil {
				return err
			}
		}

		return nil
	})
}
