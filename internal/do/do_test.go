package do_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateCluster(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		p := &do.Cluster{Name: "test", Description: "test"}
		require.NoError(t, do.CreateCluster(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		p := &do.Cluster{Name: "test", Description: "test"}
		require.NoError(t, do.CreateCluster(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
