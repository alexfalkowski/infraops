package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDocker(ctx *pulumi.Context) error {
	args := &RepositoryArgs{Topics: []string{"docker", "ruby", "golang"}}
	_, err := CreateRepository(ctx, "docker", "Common setup used for my projects.", args)

	return err
}
