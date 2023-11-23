package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createGoHealth(ctx *pulumi.Context) error {
	return createLibrary(ctx, "go-health", "Health monitoring pattern in Go.")
}

func createGoService(ctx *pulumi.Context) error {
	return createLibrary(ctx, "go-service", "A framework to build services in go.")
}
