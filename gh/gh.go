package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RepositoryArgs for gh.
type RepositoryArgs struct {
	HomepageURL string
	IsTemplate  bool
	Topics      []string
}

// CreateMasterRepository for gh.
func CreateMasterRepository(ctx *pulumi.Context, name, description string, args *RepositoryArgs) (*github.Repository, error) {
	return createRepository(ctx, name, description, "master", args)
}

// CreateMainRepository for gh.
func CreateMainRepository(ctx *pulumi.Context, name, description string, args *RepositoryArgs) (*github.Repository, error) {
	return createRepository(ctx, name, description, "main", args)
}

func createRepository(ctx *pulumi.Context, name, description, branch string, args *RepositoryArgs) (*github.Repository, error) {
	return github.NewRepository(ctx, name, &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		AllowUpdateBranch:   pulumi.Bool(true),
		DefaultBranch:       pulumi.String(branch),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String(description),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		HomepageUrl:         pulumi.String(args.HomepageURL),
		IsTemplate:          pulumi.Bool(args.IsTemplate),
		Name:                pulumi.String(name),
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
	})
}
