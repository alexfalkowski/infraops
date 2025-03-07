package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "nonnative", Description: "Allows you to keep using the power of ruby to test other systems.",
		HomepageURL: "https://alexfalkowski.github.io/nonnative",
		Checks:      gh.Checks{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
