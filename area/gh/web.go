package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "web", Description: "A website lean-thoughts.com.",
		HomepageURL: "https://alexfalkowski.github.io/web",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Template:    &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
