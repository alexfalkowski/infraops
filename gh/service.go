package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createTemplate(ctx *pulumi.Context) error {
	return createService(ctx, "go-service-template", "A template for go services.")
}

func createStatus(ctx *pulumi.Context) error {
	return createService(ctx, "status", "An alternative to https://httpstat.us/")
}
