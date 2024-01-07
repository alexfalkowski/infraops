package gh_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/gh"
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &gh.Repository{
			Name: "test", Description: "test",
			HomepageURL: "https://alexfalkowski.github.io/test", Template: gh.Template{Owner: "alexfalkowski", Repository: "go-service-template"},
		}

		r, err := gh.CreateRepository(ctx, a)
		require.NoError(t, err)
		require.NotNil(t, r)

		r.Name.ApplyT(func(n string) error {
			require.Equal(t, "test", n)

			return nil
		})

		r.Description.ApplyT(func(d *string) error {
			require.Equal(t, "test", *d)

			return nil
		})

		r.HomepageUrl.ApplyT(func(u *string) error {
			require.Equal(t, "https://alexfalkowski.github.io/test", *u)

			return nil
		})

		r.Template.ApplyT(func(tp *github.RepositoryTemplate) error {
			require.Equal(t, "alexfalkowski", tp.Owner)
			require.Equal(t, "go-service-template", tp.Repository)

			return nil
		})

		return nil
	}, pulumi.WithMocks("project", "stack", gh.Mocks(0)))

	require.NoError(t, err)
}
