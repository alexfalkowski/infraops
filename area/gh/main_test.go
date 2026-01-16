package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	config, err := gh.ReadConfiguration("gh.hjson")
	require.NoError(t, err)

	repositories := config.GetRepositories()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, repository := range repositories {
			err := gh.CreateRepository(ctx, gh.ConvertRepository(repository))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, repository := range repositories {
			err := gh.CreateRepository(ctx, gh.ConvertRepository(repository))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
