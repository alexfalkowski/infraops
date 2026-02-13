// Command apps is the Pulumi program for managing Kubernetes applications deployed by this repository.
//
// It reads `apps.hjson` from the current working directory (the Pulumi project directory) and
// provisions Kubernetes resources (for example Deployments, Services, Ingresses, and related
// supporting resources) for each application described in the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/apps`, ensuring `apps.hjson` is resolved relative to that directory.
package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := app.ReadConfiguration("apps.hjson")
		if err != nil {
			return err
		}

		for _, application := range config.GetApplications() {
			if err := app.CreateApplication(ctx, app.ConvertApplication(application)); err != nil {
				return err
			}
		}

		return nil
	})
}
