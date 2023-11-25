package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createAppConfig(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "app-config", "A place for all of my application configuration.", &RepositoryArgs{})

	return err
}
