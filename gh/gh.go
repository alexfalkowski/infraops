package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TemplateArgs for gh.
type TemplateArgs struct {
	Owner      string
	Repository string
}

func (t TemplateArgs) IsValid() bool {
	return t.Owner != "" && t.Repository != ""
}

// RepositoryArgs for gh.
type RepositoryArgs struct {
	HomepageURL string
	Template    TemplateArgs
	IsTemplate  bool
	Topics      []string
	Checks      []string
}

// CreateRepository for gh.
func CreateRepository(ctx *pulumi.Context, name, description string, args *RepositoryArgs) (*github.Repository, error) {
	return createRepository(ctx, name, description, "master", args)
}

func createRepository(ctx *pulumi.Context, name, description, branch string, args *RepositoryArgs) (*github.Repository, error) {
	r, err := repository(ctx, name, description, branch, args)
	if err != nil {
		return nil, err
	}

	if err := branchProtection(ctx, name, branch, r.NodeId, args); err != nil {
		return nil, err
	}

	return r, nil
}

func repository(ctx *pulumi.Context, name, description, branch string, args *RepositoryArgs) (*github.Repository, error) {
	a := &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		AllowUpdateBranch:   pulumi.Bool(true),
		AutoInit:            pulumi.Bool(true),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String(description),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		HomepageUrl:         pulumi.String(args.HomepageURL),
		IsTemplate:          pulumi.Bool(args.IsTemplate),
		Name:                pulumi.String(name),
		Pages: &github.RepositoryPagesArgs{
			Source: &github.RepositoryPagesSourceArgs{
				Branch: pulumi.String(branch),
			},
		},
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("enabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("enabled"),
			},
		},
		SquashMergeCommitTitle: pulumi.String("PR_TITLE"),
		Topics:                 pulumi.ToStringArray(args.Topics),
		Visibility:             pulumi.String("public"),
		VulnerabilityAlerts:    pulumi.Bool(true),
	}

	if args.Template.IsValid() {
		a.Template = github.RepositoryTemplateArgs{
			Owner:      pulumi.String(args.Template.Owner),
			Repository: pulumi.String(args.Template.Repository),
		}
	}

	return github.NewRepository(ctx, name, a)
}

func branchProtection(ctx *pulumi.Context, name, branch string, id pulumi.StringInput, args *RepositoryArgs) error {
	_, err := github.NewBranchProtection(ctx, name, &github.BranchProtectionArgs{
		Pattern:               pulumi.String(branch),
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
				Contexts: pulumi.ToStringArray(args.Checks),
				Strict:   pulumi.Bool(true),
			},
		},
	})

	return err
}
