package main

import (
	"github.com/alexfalkowski/infraops/internal/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, repository := range repositories {
			if err := gh.CreateRepository(ctx, repository); err != nil {
				return err
			}
		}

		return nil
	})
}
