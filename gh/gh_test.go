package gh_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		r, err := gh.CreateRepository(ctx, "test", "test", &gh.RepositoryArgs{})
		require.NoError(t, err)
		require.NotNil(t, r)

		return nil
	}, pulumi.WithMocks("project", "stack", gh.Mocks(0)))

	require.NoError(t, err)
}
