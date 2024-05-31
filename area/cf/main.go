package main

import (
	"github.com/alexfalkowski/infraops/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		z := &cf.Zone{
			Name:        "lean-thoughts",
			Domain:      "lean-thoughts.com",
			RecordNames: []string{"standort", "bezeichner"},
			Balancer:    "209.38.186.238",
		}

		return cf.CreateZone(ctx, z)
	})
}
