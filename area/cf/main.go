package main

import (
	"github.com/alexfalkowski/infraops/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		z := &cf.Zone{
			Name:      "lean-thoughts",
			Addresses: []string{"standort", "bezeichner"},
			Balancer:  "46.101.69.137",
		}

		return cf.CreateZone(ctx, z)
	})
}
