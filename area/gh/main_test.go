package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/internal/gh"
	test "github.com/alexfalkowski/infraops/internal/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, repository := range repositories {
			err := gh.CreateRepository(ctx, repository)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, repository := range repositories {
			err := gh.CreateRepository(ctx, repository)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
