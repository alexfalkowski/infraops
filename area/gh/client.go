package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createGoClientTemplate(ctx *pulumi.Context) error {
	checks := gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "go-client-template", Description: "A template for go clients.",
		HomepageURL: "https://alexfalkowski.github.io/go-client-template", Checks: checks,
		Template:   &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		Visibility: gh.Public, IsTemplate: true, EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createServiceControl(ctx *pulumi.Context) error {
	checks := gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "servicectl", Description: "A tool for go-service and go-service-templates.",
		HomepageURL: "https://alexfalkowski.github.io/servicectl", Checks: checks,
		Template:   &gh.Template{Owner: "alexfalkowski", Repository: "go-client-template"},
		Visibility: gh.Public, EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createKonfigControl(ctx *pulumi.Context) error {
	checks := gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "konfigctl", Description: "A tool to control https://alexfalkowski.github.io/konfig",
		HomepageURL: "https://alexfalkowski.github.io/konfigctl", Checks: checks,
		Template:   &gh.Template{Owner: "alexfalkowski", Repository: "go-client-template"},
		Visibility: gh.Public, EnablePages: true,
	}

	return gh.CreateRepository(ctx, repo)
}
