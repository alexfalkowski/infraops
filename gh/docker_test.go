package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestDocker(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createDocker(ctx)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	require.NoError(t, err)
}
