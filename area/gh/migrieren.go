package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "migrieren", Description: "Migrieren provides a way to migrate your databases.",
		HomepageURL: "https://alexfalkowski.github.io/migrieren",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
