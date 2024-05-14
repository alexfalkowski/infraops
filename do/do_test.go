package do_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/do"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		require.NoError(t, do.Configure(ctx))

		p := &do.Project{Name: "test", Description: "test"}
		require.NoError(t, do.CreateProject(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", do.Mocks(0)))

	require.NoError(t, err)
}
