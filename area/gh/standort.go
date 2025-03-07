package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "standort", Description: "Standort provides location based information.",
		HomepageURL: "https://alexfalkowski.github.io/standort",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
