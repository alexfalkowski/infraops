package main

import (
	"github.com/alexfalkowski/infraops/area/apps/lean"
	ap "github.com/alexfalkowski/infraops/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var fns ap.CreateFns

func init() {
	fns = append(fns, lean.Fns...)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}
