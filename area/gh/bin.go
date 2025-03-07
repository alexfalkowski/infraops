package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "bin", Description: "A place for common executables.",
		HomepageURL: "https://alexfalkowski.github.io/bin",
		Checks:      []string{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
