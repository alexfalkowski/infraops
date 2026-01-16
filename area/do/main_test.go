package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateCluster(t *testing.T) {
	config, err := do.ReadConfiguration("do.hjson")
	require.NoError(t, err)

	clusters := config.GetClusters()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, cluster := range clusters {
			err := do.CreateCluster(ctx, do.ConvertCluster(cluster))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, cluster := range clusters {
			err := do.CreateCluster(ctx, do.ConvertCluster(cluster))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
