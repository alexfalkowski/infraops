// Command do is the Pulumi program for managing DigitalOcean infrastructure for this repository.
//
// It reads `do.hjson` from the current working directory (the Pulumi project directory) and
// provisions DigitalOcean resources (for example VPCs and Kubernetes clusters) as described by
// the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/do`, ensuring `do.hjson` is resolved relative to that directory.
package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := do.ReadConfiguration("do.hjson")
		if err != nil {
			return err
		}

		for _, cluster := range config.GetClusters() {
			if err := do.CreateCluster(ctx, do.ConvertCluster(cluster)); err != nil {
				return err
			}
		}

		return nil
	})
}
