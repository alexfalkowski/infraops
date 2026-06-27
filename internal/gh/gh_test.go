package gh_test

import (
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

const (
	branchProtectionResourceType        = "github:index/branchProtection:BranchProtection"
	repositoryResourceType              = "github:index/repository:Repository"
	repositoryCollaboratorResourceType  = "github:index/repositoryCollaborator:RepositoryCollaborator"
	repositoryPagesResourceType         = "github:index/repositoryPages:RepositoryPages"
	repositoryVulnerabilityResourceType = "github:index/repositoryVulnerabilityAlerts:RepositoryVulnerabilityAlerts"
)

func TestCreateRepository(t *testing.T) {
	stub := &test.ResourceStub{}
	require.NoError(t, createRepository(t, fullRepository(), stub))
	requireBranchProtection(t, stub.Resources(branchProtectionResourceType))
	requirePages(t, stub.Resources(repositoryPagesResourceType))
	requireCollaborator(t, stub.Resources(repositoryCollaboratorResourceType))
}

func TestCreateRepositoryRequiresChecks(t *testing.T) {
	repository := &gh.Repository{
		Name:        "test",
		Description: "test",
		HomepageURL: "https://alexfalkowski.github.io/test",
	}

	require.Error(t, createRepository(t, repository, &test.Stub{}))
}

func TestCreateRepositoryRequiresCompleteTemplate(t *testing.T) {
	repository := &gh.Repository{
		Name:        "test",
		Description: "test",
		HomepageURL: "https://alexfalkowski.github.io/test",
		Template:    &gh.Template{Owner: "alexfalkowski"},
		Checks:      gh.Checks{"ci/circleci: build"},
	}

	require.Error(t, createRepository(t, repository, &test.Stub{}))
}

func TestCreateRepositoryReturnsRepositoryError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(repositoryResourceType)

	require.Error(t, createRepository(t, fullRepository(), stub))
	require.Len(t, stub.Resources(repositoryResourceType), 1)
}

func TestCreateRepositoryReturnsVulnerabilityAlertsError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(repositoryVulnerabilityResourceType)

	require.Error(t, createRepository(t, fullRepository(), stub))
	require.Len(t, stub.Resources(repositoryResourceType), 1)
	require.Len(t, stub.Resources(repositoryVulnerabilityResourceType), 1)
}

func TestCreateRepositoryReturnsPagesError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(repositoryPagesResourceType)

	require.Error(t, createRepository(t, fullRepository(), stub))
	require.Len(t, stub.Resources(repositoryVulnerabilityResourceType), 1)
	require.Len(t, stub.Resources(repositoryPagesResourceType), 1)
}

func TestCreateRepositoryReturnsBranchProtectionError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(branchProtectionResourceType)

	require.Error(t, createRepository(t, fullRepository(), stub))
	require.Len(t, stub.Resources(branchProtectionResourceType), 1)
}

func TestCreateRepositoryReturnsCollaboratorError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(repositoryCollaboratorResourceType)

	require.Error(t, createRepository(t, fullRepository(), stub))
	require.Len(t, stub.Resources(branchProtectionResourceType), 1)
	require.Len(t, stub.Resources(repositoryCollaboratorResourceType), 1)
}

func TestCreateRepositoryFromIncompleteTemplateConfig(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		repository := gh.ConvertRepository(&v2.Repository{
			Name:        "test",
			Description: "test",
			HomepageUrl: "https://alexfalkowski.github.io/test",
			Template:    &v2.Template{Owner: "alexfalkowski"},
			Checks:      []string{"ci/circleci: build"},
		})

		require.NotNil(t, repository.Template)
		err := gh.CreateRepository(ctx, repository)
		require.ErrorIs(t, err, gh.ErrMissingTemplate)

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)
}

func fullRepository() *gh.Repository {
	return gh.ConvertRepository(&v2.Repository{
		Name:        "test",
		Description: "test",
		HomepageUrl: "https://alexfalkowski.github.io/test",
		Visibility:  string(gh.Public),
		Collaborators: &v2.Collaborators{
			Enabled: true,
		},
		Template: &v2.Template{
			Owner:      "alexfalkowski",
			Repository: "go-service-template",
		},
		Pages: &v2.Pages{
			Enabled: true,
			Cname:   "test.alexfalkowski.dev",
		},
		Checks: []string{"ci/circleci: build"},
	})
}

func createRepository(t *testing.T, repository *gh.Repository, mocks pulumi.MockResourceMonitor) error {
	t.Helper()

	return test.RunWithMocks(func(ctx *pulumi.Context) error {
		return gh.CreateRepository(ctx, repository)
	}, mocks)
}

func requireBranchProtection(t *testing.T, protections []resource.PropertyMap) {
	t.Helper()

	require.Len(t, protections, 1)

	protection := protections[0]
	require.Equal(t, "master", test.Property(t, protection, "pattern").StringValue())
	require.True(t, test.Property(t, protection, "requiredLinearHistory").BoolValue())

	reviews := test.Property(t, protection, "requiredPullRequestReviews").ArrayValue()
	require.Len(t, reviews, 1)
	review := reviews[0].ObjectValue()
	require.True(t, test.Property(t, review, "dismissStaleReviews").BoolValue())
	require.Equal(t, 0, int(test.Property(t, review, "requiredApprovingReviewCount").NumberValue()))

	statusChecks := test.Property(t, protection, "requiredStatusChecks").ArrayValue()
	require.Len(t, statusChecks, 1)
	statusCheck := statusChecks[0].ObjectValue()
	require.True(t, test.Property(t, statusCheck, "strict").BoolValue())
	require.Equal(t, []string{"ci/circleci: build"}, test.StringValues(test.Property(t, statusCheck, "contexts").ArrayValue()))
}

func requirePages(t *testing.T, pages []resource.PropertyMap) {
	t.Helper()

	require.Len(t, pages, 1)

	page := pages[0]
	require.Equal(t, "legacy", test.Property(t, page, "buildType").StringValue())
	require.Equal(t, "test.alexfalkowski.dev", test.Property(t, page, "cname").StringValue())
	require.Equal(t, "test", test.Property(t, page, "repository").StringValue())

	source := test.Property(t, page, "source").ObjectValue()
	require.Equal(t, "master", test.Property(t, source, "branch").StringValue())
}

func requireCollaborator(t *testing.T, collaborators []resource.PropertyMap) {
	t.Helper()

	require.Len(t, collaborators, 1)

	collaborator := collaborators[0]
	require.Equal(t, "push", test.Property(t, collaborator, "permission").StringValue())
	require.Equal(t, "alexfalkowski/test", test.Property(t, collaborator, "repository").StringValue())
	require.Equal(t, "lean-thoughts-ci", test.Property(t, collaborator, "username").StringValue())
}
