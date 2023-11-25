package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNonnative(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "nonnative", "Allows you to keep using the power of ruby to test other systems.", &RepositoryArgs{})

	return err
}
