package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "alexfalkowski.github.io", Description: "A site for my profile.",
		HomepageURL: "https://alexfalkowski.github.io",
		Checks:      gh.Checks{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
