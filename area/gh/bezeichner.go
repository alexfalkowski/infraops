package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "bezeichner", Description: "Bezeichner takes care of identifiers used in your services.",
		HomepageURL: "https://alexfalkowski.github.io/bezeichner",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
