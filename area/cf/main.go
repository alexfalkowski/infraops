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
