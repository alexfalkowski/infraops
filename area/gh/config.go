package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "app-config", Description: "A place for all of my application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/app-config",
		Checks:      []string{"ci/circleci: verify-config"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
