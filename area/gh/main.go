package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(run("gh.hjson"))
}

func run(path string) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		config, err := gh.ReadConfiguration(path)
		if err != nil {
			return err
		}

		repositories := make([]*gh.Repository, len(config.GetRepositories()))
		for i, repository := range config.GetRepositories() {
			repositories[i] = gh.ConvertRepository(repository)
		}

		return gh.CreateRepositories(ctx, repositories)
	}
}
