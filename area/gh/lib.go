package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNonnative(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "nonnative", Description: "Allows you to keep using the power of ruby to test other systems.",
		HomepageURL: "https://alexfalkowski.github.io/nonnative", Checks: []string{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createGoHealth(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "go-health", Description: "Health monitoring pattern in Go.",
		HomepageURL: "https://alexfalkowski.github.io/go-health", Checks: []string{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createGoService(ctx *pulumi.Context) error {
	repo := &gh.Repository{
		Name: "go-service", Description: "A framework to build services in go.",
		HomepageURL: "https://alexfalkowski.github.io/go-service", Checks: []string{"ci/circleci: build"},
		Visibility: "public", EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}
