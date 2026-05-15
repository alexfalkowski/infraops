package gh_test

import (
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/alexfalkowski/infraops/v2/internal/test"
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
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test", HomepageURL: "https://alexfalkowski.github.io/test",
		}

		err := gh.CreateRepository(ctx, a)
		require.Error(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
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
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
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
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}

func TestCreateRepositoryFromIncompleteTemplateConfig(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := gh.ConvertRepository(&v2.Repository{
			Name:        "test",
			Description: "test",
			HomepageUrl: "https://alexfalkowski.github.io/test",
			Template:    &v2.Template{Owner: "alexfalkowski"},
			Checks:      []string{"ci/circleci: build"},
		})

		require.NotNil(t, a.Template)
		err := gh.CreateRepository(ctx, a)
		require.ErrorIs(t, err, gh.ErrMissingTemplate)

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)
}
