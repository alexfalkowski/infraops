package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/test"
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
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, fn := range fns {
			err := fn(ctx)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
