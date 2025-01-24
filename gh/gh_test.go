package gh_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/gh"
	"github.com/alexfalkowski/infraops/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
			Template: &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
			Checks:   gh.Checks{"ci/circleci: build"},
		}

		err := gh.CreateRepository(ctx, a)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.Mocks(0)))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
		}

		err := gh.CreateRepository(ctx, a)
		require.Error(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.Mocks(0)))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
			Template: &gh.Template{Owner: "alexfalkowski"},
			Checks:   gh.Checks{"ci/circleci: build"},
		}

		err := gh.CreateRepository(ctx, a)
		require.Error(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.Mocks(0)))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
			Template: &gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
			Checks:   gh.Checks{"ci/circleci: build"},
		}

		err := gh.CreateRepository(ctx, a)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.BadMocks(0)))
	require.Error(t, err)
}
