package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createTemplate(ctx *pulumi.Context) error {
	return createService(ctx, "go-service-template", "A template for go services.", true)
}

func createStatus(ctx *pulumi.Context) error {
	return createService(ctx, "status", "An alternative to https://httpstat.us/", false)
}

func createStandort(ctx *pulumi.Context) error {
	return createService(ctx, "standort", "Standort provides location based information.", false)
}

func createAuth(ctx *pulumi.Context) error {
	return createService(ctx, "auth", "Auth provides all your authn and authz needs.", false)
}

func createKonfig(ctx *pulumi.Context) error {
	return createService(ctx, "konfig", "Konfig is a configuration system for application configuration.", false)
}

func createMigrieren(ctx *pulumi.Context) error {
	return createService(ctx, "migrieren", "Migrieren provides a way to migrate your databases.", false)
}
