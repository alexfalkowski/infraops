package main

import "github.com/alexfalkowski/infraops/internal/gh"

func init() {
	RegisterRepository(&gh.Repository{
		Name: "konfigctl", Description: "A tool to control https://alexfalkowski.github.io/konfig.",
		HomepageURL: "https://alexfalkowski.github.io/konfigctl",
		Checks:      gh.Checks{"ci/circleci: build-client", "ci/circleci: build-docker"},
		Template:    &gh.Template{Owner: "alexfalkowski", Repository: "go-client-template"},
		Visibility:  gh.Public,
		EnablePages: true,
	})

	RegisterRepository(&gh.Repository{
		Name: "konfig", Description: "Konfig is a configuration system for application configuration.",
		HomepageURL: "https://alexfalkowski.github.io/konfig",
		Checks:      gh.Checks{"ci/circleci: build-service", "ci/circleci: build-docker", "ci/circleci: features-grpc", "ci/circleci: features-http", "ci/circleci: features-coverage"},
		Visibility:  gh.Public,
		EnablePages: true,
	})
}
