// Command cf is the Pulumi program for managing Cloudflare infrastructure for this repository.
//
// It reads `cf.hjson` from the current working directory (the Pulumi project directory) and
// provisions Cloudflare resources (for example zones, DNS records, and R2 buckets/custom domains)
// as described by the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/cf`, ensuring `cf.hjson` is resolved relative to that directory.
package main

import (
	"github.com/alexfalkowski/infraops/v2/internal/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config, err := cf.ReadConfiguration("cf.hjson")
		if err != nil {
			return err
		}

		for _, zone := range config.GetBalancerZones() {
			if err := cf.CreateBalancerZone(ctx, cf.ConvertBalancerZone(zone)); err != nil {
				return err
			}
		}

		for _, zone := range config.GetPageZones() {
			if err := cf.CreatePageZone(ctx, cf.ConvertPageZone(zone)); err != nil {
				return err
			}
		}

		for _, bucket := range config.GetBuckets() {
			if err := cf.CreateBucket(ctx, cf.ConvertBucket(bucket)); err != nil {
				return err
			}
		}

		return nil
	})
}
