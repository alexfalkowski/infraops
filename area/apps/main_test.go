package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestFns(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, fn := range fns {
			err := fn(ctx)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", app.Mocks(0)))

	require.NoError(t, err)
}
