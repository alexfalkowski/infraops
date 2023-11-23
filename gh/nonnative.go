package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNonnative(ctx *pulumi.Context) error {
	return createLibrary(ctx, "nonnative", "Allows you to keep using the power of ruby to test other systems.")
}
