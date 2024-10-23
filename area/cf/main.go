package main

import (
	"github.com/alexfalkowski/infraops/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		zones := []*cf.Zone{
			{
				Name:        "lean-thoughts",
				Domain:      "lean-thoughts.com",
				RecordNames: []string{"standort", "bezeichner", "web"},
				Balancer:    "209.38.186.238",
			},
			{
				Name:        "sasha-adventures",
				Domain:      "sasha-adventures.com",
				RecordNames: []string{"www"},
				Balancer:    "209.38.186.238",
			},
			{
				Name:        "afalkowski",
				Domain:      "afalkowski.com",
				RecordNames: []string{"www"},
				Balancer:    "209.38.186.238",
			},
		}

		for _, z := range zones {
			if err := cf.CreateZone(ctx, z); err != nil {
				return err
			}
		}

		return nil
	})
}
