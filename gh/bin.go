package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createBin(ctx *pulumi.Context) error {
	args := &RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/bin"}
	_, err := CreateRepository(ctx, "bin", "A place for common executables.", args)

	return err
}
