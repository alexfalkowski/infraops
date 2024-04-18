package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

var fns = []createFn{
	createTest,
	createPages, createInfraOps,
	createDocker, createAppConfig,
	createBin, createNonnative,
	createGoHealth, createGoService,
	createGoServiceTemplate, createStatus,
	createStandort, createAuth,
	createKonfig, createMigrieren,
}

func createTest(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "test", Description: "Test repository.",
		HomepageURL: "https://alexfalkowski.github.io/test", Checks: checks,
		Template: gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createInfraOps(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "infraops", Description: "A place where all infrastructure is taken care of.",
		HomepageURL: "https://alexfalkowski.github.io/infraops", Checks: []string{"ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createPages(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "alexfalkowski.github.io", Description: "A site for my profile.",
		HomepageURL: "https://alexfalkowski.github.io",
	}

	return gh.CreateRepository(ctx, repo)
}

func createDocker(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "docker", Description: "Common setup used for my projects.",
		Topics: []string{"docker", "ruby", "golang"}, Checks: []string{"ci/circleci: lint", "ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createAppConfig(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "app-config", Description: "A place for all of my application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/app-config",
	}

	return gh.CreateRepository(ctx, repo)
}

func createBin(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "bin", Description: "A place for common executables.",
		HomepageURL: "https://alexfalkowski.github.io/bin", Checks: []string{"ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createNonnative(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "nonnative", Description: "Allows you to keep using the power of ruby to test other systems.",
		HomepageURL: "https://alexfalkowski.github.io/nonnative", Checks: []string{"ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createGoHealth(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "go-health", Description: "Health monitoring pattern in Go.",
		HomepageURL: "https://alexfalkowski.github.io/go-health", Checks: []string{"ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createGoService(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "go-service", Description: "A framework to build services in go.",
		HomepageURL: "https://alexfalkowski.github.io/go-service", Checks: []string{"ci/circleci: build"},
	}

	return gh.CreateRepository(ctx, repo)
}

func createGoServiceTemplate(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "go-service-template", Description: "A template for go services.",
		HomepageURL: "https://alexfalkowski.github.io/go-service-template", Checks: checks, IsTemplate: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createStatus(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "status", Description: "An alternative to https://httpstat.us/",
		HomepageURL: "https://alexfalkowski.github.io/status", Checks: checks,
	}

	return gh.CreateRepository(ctx, repo)
}

func createStandort(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coveralls"}
	repo := &gh.Repository{
		Name: "standort", Description: "Standort provides location based information.",
		HomepageURL: "https://alexfalkowski.github.io/standort", Checks: checks,
	}

	return gh.CreateRepository(ctx, repo)
}

func createAuth(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "auth", Description: "Auth provides all your authn and authz needs.",
		HomepageURL: "https://alexfalkowski.github.io/auth", Checks: checks,
	}

	return gh.CreateRepository(ctx, repo)
}

func createKonfig(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coveralls"}
	repo := &gh.Repository{
		Name: "konfig", Description: "Konfig is a configuration system for application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/konfig", Checks: checks,
	}

	return gh.CreateRepository(ctx, repo)
}

func createMigrieren(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "migrieren", Description: "Migrieren provides a way to migrate your databases.",
		HomepageURL: "https://alexfalkowski.github.io/migrieren", Checks: checks,
	}

	return gh.CreateRepository(ctx, repo)
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
