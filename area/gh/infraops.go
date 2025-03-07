package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "infraops", Description: "A place where all infrastructure is taken care of.",
		HomepageURL: "https://alexfalkowski.github.io/infraops",
		Checks:      gh.Checks{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
