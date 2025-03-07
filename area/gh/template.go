package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "go-service-template", Description: "A template for go services.",
		HomepageURL: "https://alexfalkowski.github.io/go-service-template",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Visibility:  gh.Public,
		IsTemplate:  true,
		EnablePages: true,
	})

	RegisterRepository(&gh.Repository{
		Name: "go-client-template", Description: "A template for go clients.",
		HomepageURL: "https://alexfalkowski.github.io/go-client-template",
		Checks:      gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"},
		Template:    &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		Visibility:  gh.Public,
		IsTemplate:  true,
		EnablePages: true,
	})
}
