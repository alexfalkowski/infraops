package gh

import (
	"errors"

	errs "github.com/alexfalkowski/infraops/internal/errors"
	"github.com/alexfalkowski/infraops/internal/runtime"
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const master = "master"

var (
	// ErrMissingTemplate for gh.
	ErrMissingTemplate = errors.New("missing template")

	// ErrMissingChecks for gh.
	ErrMissingChecks = errors.New("missing checks")

	// Public visibility.
	Public = Visibility("public")

	// Private visibility.
	Private = Visibility("private")
)

type (
	// Visibility of the repositories.
	Visibility string

	// Template for gh.
	Template struct {
		Owner      string
		Repository string
	}

	// Checks for gh.
	Checks []string

	// Repository for gh.
	Repository struct {
		Template    *Template
		Name        string
		Description string
		HomepageURL string
		Visibility  Visibility
		Topics      []string
		Checks      Checks
		IsTemplate  bool
		EnablePages bool
		Archived    bool
	}
)

// Valid if no error is returned.
func (t *Template) Valid() error {
	if t.Owner == "" || t.Repository == "" {
		return ErrMissingTemplate
	}

	return nil
}

// Valid if no error is returned.
func (c Checks) Valid() error {
	if len(c) == 0 {
		return ErrMissingChecks
	}

	return nil
}

// CreateRepository for gh.
func CreateRepository(ctx *pulumi.Context, repo *Repository) error {
	r, err := repository(ctx, repo)
	if err != nil {
		return errs.Prefix(repo.Name, err)
	}

	return errs.Prefix(repo.Name, branchProtection(ctx, r.NodeId, repo))
}

//nolint:nonamedreturns
func repository(ctx *pulumi.Context, repo *Repository) (repository *github.Repository, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	t, err := template(repo)
	runtime.Must(err)

	err = repo.Checks.Valid()
	runtime.Must(err)

	args := &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		AllowUpdateBranch:   pulumi.Bool(true),
		AutoInit:            pulumi.Bool(true),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String(repo.Description),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		HomepageUrl:         pulumi.String(repo.HomepageURL),
		IsTemplate:          pulumi.Bool(repo.IsTemplate),
		Name:                pulumi.String(repo.Name),
		Archived:            pulumi.Bool(repo.Archived),
		Pages:               pages(repo),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("enabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("enabled"),
			},
		},
		SquashMergeCommitTitle: pulumi.String("PR_TITLE"),
		Template:               t,
		Topics:                 pulumi.ToStringArray(repo.Topics),
		Visibility:             pulumi.String(repo.Visibility),
		VulnerabilityAlerts:    pulumi.Bool(true),
	}

	repository, err = github.NewRepository(ctx, repo.Name, args)
	runtime.Must(err)

	return //nolint:nakedret
}

func branchProtection(ctx *pulumi.Context, id pulumi.StringInput, repo *Repository) error {
	args := &github.BranchProtectionArgs{
		Pattern:               pulumi.String(master),
		RepositoryId:          id,
		RequiredLinearHistory: pulumi.Bool(true),
		RequiredPullRequestReviews: github.BranchProtectionRequiredPullRequestReviewArray{
			&github.BranchProtectionRequiredPullRequestReviewArgs{
				DismissStaleReviews:          pulumi.Bool(true),
				RequiredApprovingReviewCount: pulumi.Int(0),
			},
		},
		RequiredStatusChecks: github.BranchProtectionRequiredStatusCheckArray{
			&github.BranchProtectionRequiredStatusCheckArgs{
				Contexts: pulumi.ToStringArray(repo.Checks),
				Strict:   pulumi.Bool(true),
			},
		},
	}
	_, err := github.NewBranchProtection(ctx, repo.Name, args)

	return err
}

//nolint:nilnil
func template(repo *Repository) (*github.RepositoryTemplateArgs, error) {
	if repo.Template == nil {
		return nil, nil
	}

	if err := repo.Template.Valid(); err != nil {
		return nil, err
	}

	args := &github.RepositoryTemplateArgs{
		Owner:      pulumi.String(repo.Template.Owner),
		Repository: pulumi.String(repo.Template.Repository),
	}

	return args, nil
}

func pages(repo *Repository) *github.RepositoryPagesArgs {
	if !repo.EnablePages {
		return nil
	}

	return &github.RepositoryPagesArgs{
		Source: &github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(master),
		},
	}
}
