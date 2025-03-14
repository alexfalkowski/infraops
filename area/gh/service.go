package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "go-service", Description: "A framework to build services in go.",
		HomepageURL: "https://alexfalkowski.github.io/go-service",
		Checks:      gh.Checks{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
