package gh

import (
	"errors"
	"fmt"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/config"
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

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Github, error) {
	var configuration v2.Github
	err := config.Read(path, &configuration)

	return &configuration, err
}

// ConvertRepository converts a v2.Repository to a Repository.
func ConvertRepository(r *v2.Repository) *Repository {
	repository := &Repository{
		Name:        r.GetName(),
		Description: r.GetDescription(),
		HomepageURL: r.GetHomepageUrl(),
		Visibility:  Visibility(r.GetVisibility()),
		Topics:      r.GetTopics(),
		Checks:      Checks(r.GetChecks()),
		IsTemplate:  r.GetIsTemplate(),
		EnablePages: r.GetEnablePages(),
		Archived:    r.GetArchived(),
	}

	if template := r.GetTemplate(); template != nil {
		owner := template.GetOwner()
		repo := template.GetRepository()

		if owner != "" && repo != "" {
			repository.Template = &Template{
				Owner:      owner,
				Repository: repo,
			}
		}
	}

	return repository
}

// CreateRepository for gh.
func CreateRepository(ctx *pulumi.Context, repo *Repository) error {
	r, err := repository(ctx, repo)
	if err != nil {
		return fmt.Errorf("%v: %w", repo.Name, err)
	}

	if err := branchProtection(ctx, r.NodeId, repo); err != nil {
		return fmt.Errorf("%v: %w", repo.Name, err)
	}

	return nil
}

func repository(ctx *pulumi.Context, repo *Repository) (*github.Repository, error) {
	t, err := template(repo)
	if err != nil {
		return nil, err
	}

	if err := repo.Checks.Valid(); err != nil {
		return nil, err
	}

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

	return github.NewRepository(ctx, repo.Name, args)
}

func branchProtection(ctx *pulumi.Context, id pulumi.StringInput, repo *Repository) error {
	args := &github.BranchProtectionArgs{
		EnforceAdmins:         pulumi.Bool(true),
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
