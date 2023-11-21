package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		args := &github.RepositoryArgs{
			Description: pulumi.String("Demo Repository for GitHub"),
		}

		repository, err := github.NewRepository(ctx, "demo-repo", args)
		if err != nil {
			return err
		}

		ctx.Export("repository", repository.Name)

		return nil
	})
}
