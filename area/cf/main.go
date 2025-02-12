package main

import (
	"github.com/alexfalkowski/infraops/internal/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		lean := &cf.BalancerZone{
			Name:        "lean-thoughts",
			Domain:      "lean-thoughts.com",
			RecordNames: []string{"standort", "bezeichner", "web"},
			IP:          "209.38.186.238",
		}

		if err := cf.CreateBalancerZone(ctx, lean); err != nil {
			return err
		}

		pages := []*cf.PageZone{
			{
				Name:   "sasha-adventures",
				Domain: "sasha-adventures.com",
				Host:   "sasha-adventures.github.io",
			},
			{
				Name:   "afalkowski",
				Domain: "afalkowski.com",
				Host:   "alexfalkowski.github.io",
			},
		}

		for _, p := range pages {
			if err := cf.CreatePageZone(ctx, p); err != nil {
				return err
			}
		}

		return nil
	})
}
