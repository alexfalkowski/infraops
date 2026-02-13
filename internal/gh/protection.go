package gh

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// branchProtection applies branch protection rules to the default branch for repo.
//
// It enforces:
//   - linear history
//   - pull request reviews (with stale reviews dismissed)
//   - required status checks as specified by repo.Checks (strict mode enabled)
//
// The branch name is taken from the package-level master constant.
func branchProtection(ctx *pulumi.Context, id pulumi.StringInput, repo *Repository) error {
	args := &github.BranchProtectionArgs{
		Pattern:               pulumi.String(master),
		RepositoryId:          id,
		RequiredLinearHistory: inputs.Yes,
		RequiredPullRequestReviews: github.BranchProtectionRequiredPullRequestReviewArray{
			&github.BranchProtectionRequiredPullRequestReviewArgs{
				DismissStaleReviews:          inputs.Yes,
				RequiredApprovingReviewCount: pulumi.Int(0),
			},
		},
		RequiredStatusChecks: github.BranchProtectionRequiredStatusCheckArray{
			&github.BranchProtectionRequiredStatusCheckArgs{
				Contexts: pulumi.ToStringArray(repo.Checks),
				Strict:   inputs.Yes,
			},
		},
	}
	_, err := github.NewBranchProtection(ctx, repo.Name, args)

	return err
}
