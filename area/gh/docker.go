package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "docker", Description: "Common setup used for my projects.",
		Topics:      []string{"docker", "ruby", "golang"},
		Checks:      gh.Checks{"ci/circleci: lint", "ci/circleci: build"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
