package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type createFn func(ctx *pulumi.Context) error

var fns = []createFn{
	createPages, createInfraOps,
	createDocker, createAppConfig,
	createBin, createNonnative,
	createGoHealth, createGoService,
	createGoServiceTemplate, createStatus,
	createStandort, createAuth,
	createKonfig, createMigrieren,
}

func createInfraOps(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/infraops", Checks: []string{"ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "infraops", "A place where all infrastructure is taken care of.", args)

	return err
}

func createPages(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io"}
	_, err := gh.CreateRepository(ctx, "alexfalkowski.github.io", "A site for my profile.", args)

	return err
}

func createDocker(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{Topics: []string{"docker", "ruby", "golang"}, Checks: []string{"ci/circleci: lint", "ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "docker", "Common setup used for my projects.", args)

	return err
}

func createAppConfig(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/app-config"}
	_, err := gh.CreateRepository(ctx, "app-config", "A place for all of my application configuration.", args)

	return err
}

func createBin(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/bin", Checks: []string{"ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "bin", "A place for common executables.", args)

	return err
}

func createNonnative(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/nonnative", Checks: []string{"ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "nonnative", "Allows you to keep using the power of ruby to test other systems.", args)

	return err
}

func createGoHealth(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/go-health", Checks: []string{"ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "go-health", "Health monitoring pattern in Go.", args)

	return err
}

func createGoService(ctx *pulumi.Context) error {
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/go-service", Checks: []string{"ci/circleci: build"}}
	_, err := gh.CreateRepository(ctx, "go-service", "A framework to build services in go.", args)

	return err
}

func createGoServiceTemplate(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/go-service-template", Checks: checks, IsTemplate: true}
	_, err := gh.CreateRepository(ctx, "go-service-template", "A template for go services.", args)

	return err
}

func createStatus(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/status", Checks: checks}
	_, err := gh.CreateRepository(ctx, "status", "An alternative to https://httpstat.us/", args)

	return err
}

func createStandort(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coveralls"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/standort", Checks: checks}
	_, err := gh.CreateRepository(ctx, "standort", "Standort provides location based information.", args)

	return err
}

func createAuth(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/auth", Checks: checks}
	_, err := gh.CreateRepository(ctx, "auth", "Auth provides all your authn and authz needs.", args)

	return err
}

func createKonfig(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coveralls"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/konfig", Checks: checks}
	_, err := gh.CreateRepository(ctx, "konfig", "Konfig is a configuration system for application configuration.", args)

	return err
}

func createMigrieren(ctx *pulumi.Context) error {
	checks := []string{"ci/circleci: build-service", "ci/circleci: build-docker"}
	args := &gh.RepositoryArgs{HomepageURL: "https://alexfalkowski.github.io/migrieren", Checks: checks}
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
