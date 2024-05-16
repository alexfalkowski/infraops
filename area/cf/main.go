package main

import (
	"github.com/alexfalkowski/infraops/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		z := &cf.Zone{
			Name:      "lean-thoughts",
			Addresses: []string{"api.standort", "grpc.standort"},
			Balancer:  "138.68.124.205",
		}

		return cf.CreateZone(ctx, z)
	})
}
