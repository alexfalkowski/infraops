package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createGoServiceTemplate(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "go-service-template", Description: "A template for go services.",
		HomepageURL: "https://alexfalkowski.github.io/go-service-template", Checks: checks,
		Visibility: "public", IsTemplate: true, EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createStatus(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "status", Description: "An alternative to https://httpstat.us/",
		HomepageURL: "https://alexfalkowski.github.io/status", Checks: checks,
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createStandort(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "standort", Description: "Standort provides location based information.",
		HomepageURL: "https://alexfalkowski.github.io/standort", Checks: checks,
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createAuth(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "auth", Description: "Auth provides all your authn and authz needs.",
		HomepageURL: "https://alexfalkowski.github.io/auth", Checks: checks,
		Visibility: "public", EnablePages: true, Archived: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createKonfig(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coverage"}
	repo := &gh.Repository{
		Name: "konfig", Description: "Konfig is a configuration system for application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/konfig", Checks: checks,
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createMigrieren(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "migrieren", Description: "Migrieren provides a way to migrate your databases.",
		HomepageURL: "https://alexfalkowski.github.io/migrieren", Checks: checks,
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createBezeichner(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "bezeichner", Description: "Bezeichner takes care of identifiers used in your services.",
		HomepageURL: "https://alexfalkowski.github.io/bezeichner", Checks: checks,
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}
