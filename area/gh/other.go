package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createInfraOps(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "infraops", Description: "A place where all infrastructure is taken care of.",
		HomepageURL: "https://alexfalkowski.github.io/infraops", Checks: gh.Checks{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createSite(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "alexfalkowski.github.io", Description: "A site for my profile.",
		HomepageURL: "https://alexfalkowski.github.io", Checks: gh.Checks{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createDocker(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "docker", Description: "Common setup used for my projects.",
		Topics: []string{"docker", "ruby", "golang"}, Checks: gh.Checks{"ci/circleci: lint", "ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createAppConfig(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "app-config", Description: "A place for all of my application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/app-config", Checks: []string{"ci/circleci: verify-config"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createBin(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "bin", Description: "A place for common executables.",
		HomepageURL: "https://alexfalkowski.github.io/bin", Checks: []string{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}
