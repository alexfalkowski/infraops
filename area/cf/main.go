package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		a, err := cloudflare.NewAccount(ctx, "main", &cloudflare.AccountArgs{
			Name: pulumi.String("Alexrfalkowski@gmail.com's Account"),
		})
		if err != nil {
			return err
		}

		_, err = cloudflare.NewZone(ctx, "lean-thoughts", &cloudflare.ZoneArgs{
			AccountId: a.ID(),
			Plan:      pulumi.String("free"),
			Zone:      pulumi.String("lean-thoughts.com"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
