package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createTemplate(ctx *pulumi.Context) error {
	args := &RepositoryArgs{IsTemplate: true}
	_, err := CreateRepository(ctx, "go-service-template", "A template for go services.", args)

	return err
}

func createStatus(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "status", "An alternative to https://httpstat.us/", &RepositoryArgs{})

	return err
}

func createStandort(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "standort", "Standort provides location based information.", &RepositoryArgs{})

	return err
}

func createAuth(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "auth", "Auth provides all your authn and authz needs.", &RepositoryArgs{})

	return err
}

func createKonfig(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "konfig", "Konfig is a configuration system for application configuration.", &RepositoryArgs{})

	return err
}

func createMigrieren(ctx *pulumi.Context) error {
	_, err := CreateRepository(ctx, "migrieren", "Migrieren provides a way to migrate your databases.", &RepositoryArgs{})

	return err
}
