package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const master = "master"

// Template for gh.
type Template struct {
	Owner      string
	Repository string
}

func (t Template) IsValid() bool {
	return t.Owner != "" && t.Repository != ""
}

// Repository for gh.
type Repository struct {
	Name        string
	Description string
	HomepageURL string
	Template    Template
	IsTemplate  bool
	Topics      []string
	Checks      []string
	EnablePages bool
}

// CreateRepository for gh.
func CreateRepository(ctx *pulumi.Context, repo *Repository) error {
	r, err := repository(ctx, repo)
	if err != nil {
		return err
	}

	return branchProtection(ctx, r.NodeId, repo)
}

func repository(ctx *pulumi.Context, repo *Repository) (*github.Repository, error) {
	a := &github.RepositoryArgs{
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
		Template:               template(repo),
		Topics:                 pulumi.ToStringArray(repo.Topics),
		Visibility:             pulumi.String("public"),
		VulnerabilityAlerts:    pulumi.Bool(true),
	}

	return github.NewRepository(ctx, repo.Name, a)
}

func branchProtection(ctx *pulumi.Context, id pulumi.StringInput, repo *Repository) error {
	_, err := github.NewBranchProtection(ctx, repo.Name, &github.BranchProtectionArgs{
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
	})

	return err
}

func template(repo *Repository) *github.RepositoryTemplateArgs {
	if !repo.Template.IsValid() {
		return nil
	}

	return &github.RepositoryTemplateArgs{Owner: pulumi.String(repo.Template.Owner), Repository: pulumi.String(repo.Template.Repository)}
}

func pages(repo *Repository) *github.RepositoryPagesArgs {
	if !repo.EnablePages {
		return nil
	}

	return &github.RepositoryPagesArgs{Source: &github.RepositoryPagesSourceArgs{Branch: pulumi.String(master)}}
}
