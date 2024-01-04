package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

var fns = []createFn{
	createDocker, createAppConfig,
	createBin, createNonnative,
	createGoHealth, createGoService,
	createGoServiceTemplate, createStatus,
	createStandort, createAuth,
	createKonfig, createMigrieren,
}

func createDocker(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{Topics: []string{"docker", "ruby", "golang"}}
	_, err := gh.CreateRepository(ctx, "docker", "Common setup used for my projects.", args)

	return err
}

func createAppConfig(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/app-config"}
	_, err := gh.CreateRepository(ctx, "app-config", "A place for all of my application configuration.", args)

	return err
}

func createBin(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/bin"}
	_, err := gh.CreateRepository(ctx, "bin", "A place for common executables.", args)

	return err
}

func createNonnative(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/nonnative"}
	_, err := gh.CreateRepository(ctx, "nonnative", "Allows you to keep using the power of ruby to test other systems.", args)

	return err
}

func createGoHealth(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/go-health"}
	_, err := gh.CreateRepository(ctx, "go-health", "Health monitoring pattern in Go.", args)

	return err
}

func createGoService(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/go-service"}
	_, err := gh.CreateRepository(ctx, "go-service", "A framework to build services in go.", args)

	return err
}

func createGoServiceTemplate(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/go-service-template", IsTemplate: true}
	_, err := gh.CreateRepository(ctx, "go-service-template", "A template for go services.", args)

	return err
}

func createStatus(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/status"}
	_, err := gh.CreateRepository(ctx, "status", "An alternative to https://httpstat.us/", args)

	return err
}

func createStandort(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/standort"}
	_, err := gh.CreateRepository(ctx, "standort", "Standort provides location based information.", args)

	return err
}

func createAuth(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/auth"}
	_, err := gh.CreateRepository(ctx, "auth", "Auth provides all your authn and authz needs.", args)

	return err
}

func createKonfig(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/konfig"}
	_, err := gh.CreateRepository(ctx, "konfig", "Konfig is a configuration system for application configuration.", args)

	return err
}

func createMigrieren(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://github.com/alexfalkowski/migrieren"}
	_, err := gh.CreateRepository(ctx, "migrieren", "Migrieren provides a way to migrate your databases.", args)

	return err
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
