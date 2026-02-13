// Command gh is the Pulumi program for managing GitHub infrastructure for this repository.
//
// It reads `gh.hjson` from the current working directory (the Pulumi project directory)
// and provisions GitHub repositories and their associated configuration (for example branch
// protection, collaborators, and Pages) as described by the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/gh`, ensuring `gh.hjson` is resolved relative to that directory.
package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := gh.ReadConfiguration("gh.hjson")
		if err != nil {
			return err
		}

		for _, repository := range config.GetRepositories() {
			if err := gh.CreateRepository(ctx, gh.ConvertRepository(repository)); err != nil {
				return err
			}
		}

		return nil
	})
}
