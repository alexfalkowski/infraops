package main

import (
	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createIDPControl(ctx *pulumi.Context) error {
	checks := gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "idpctl", Description: "A tool to control https://alexfalkowski.github.io/idp",
		HomepageURL: "https://alexfalkowski.github.io/idpctl", Checks: checks,
		Template:   &gh.Template{Owner: "alexfalkowski", Repository: "go-client-template"},
		Visibility: gh.Public, EnablePages: false, Archived: true,
	}

	return gh.CreateRepository(ctx, repo)
}

func createIDPService(ctx *pulumi.Context) error {
	checks := gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"}
	repo := &gh.Repository{
		Name: "idpd", Description: "Internal Developer Platform.",
		HomepageURL: "https://alexfalkowski.github.io/idpd", Checks: checks,
		Template:   &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		Visibility: gh.Public, EnablePages: false, Archived: true,
	}

	return gh.CreateRepository(ctx, repo)
}
