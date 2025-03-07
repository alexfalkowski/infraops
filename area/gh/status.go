package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "status", Description: "An alternative to https://httpstat.us/.",
		HomepageURL: "https://alexfalkowski.github.io/status",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
