package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func branchProtection(ctx *pulumi.Context, id pulumi.StringInput, repo *Repository) error {
	args := &github.BranchProtectionArgs{
		EnforceAdmins:         pulumi.Bool(false),
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
