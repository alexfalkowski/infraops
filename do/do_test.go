package do_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/do"
	"github.com/alexfalkowski/infraops/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		require.NoError(t, do.Configure(ctx))

		p := &do.Project{Name: "test", Description: "test"}
		require.NoError(t, do.CreateProject(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		require.NoError(t, do.Configure(ctx))

		p := &do.Project{Name: "test", Description: "test"}
		require.NoError(t, do.CreateProject(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
