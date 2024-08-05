package main

import (
	ap "github.com/alexfalkowski/infraops/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var fns ap.CreateFns

func init() {
	fns = append(fns, Fns...)
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
