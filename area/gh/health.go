package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "go-health", Description: "Health monitoring pattern in go.",
		HomepageURL: "https://alexfalkowski.github.io/go-health",
		Checks:      gh.Checks{"ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
