package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := do.ReadConfiguration("do.yaml")
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
