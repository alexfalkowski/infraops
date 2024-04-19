package gh_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
			Template: gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		}

		err := gh.CreateRepository(ctx, a)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", gh.Mocks(0)))

	require.NoError(t, err)
}
