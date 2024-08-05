package pulumi

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type (
	// CreateFn for pulumi.
	CreateFn func(ctx *pulumi.Context) error

	// CreateFns for pulumi.
	CreateFns []CreateFn
)
