package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createGoHealth(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "go-health", "Health monitoring pattern in Go.", &RepositoryArgs{})

	return err
}

func createGoService(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "go-service", "A framework to build services in go.", &RepositoryArgs{})

	return err
}
