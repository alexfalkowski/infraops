package gh

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func repository(ctx *pulumi.Context, repo *Repository) (*github.Repository, error) {
	t, err := template(repo)
	if err != nil {
		return nil, err
	}

	if err := repo.Checks.Valid(); err != nil {
		return nil, err
	}

	args := &github.RepositoryArgs{
		AllowMergeCommit:    inputs.No,
		AllowRebaseMerge:    inputs.No,
		AllowUpdateBranch:   inputs.Yes,
		AutoInit:            inputs.Yes,
		AllowAutoMerge:      inputs.Yes,
		DeleteBranchOnMerge: inputs.Yes,
		Description:         pulumi.String(repo.Description),
		HasDownloads:        inputs.Yes,
		HasIssues:           inputs.Yes,
		HasProjects:         inputs.Yes,
		HasWiki:             inputs.Yes,
		HomepageUrl:         pulumi.String(repo.HomepageURL),
		IsTemplate:          pulumi.Bool(repo.IsTemplate),
		Name:                pulumi.String(repo.Name),
		Archived:            pulumi.Bool(repo.Archived),
		Pages:               pages(repo),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: inputs.Enabled,
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: inputs.Enabled,
			},
		},
		SquashMergeCommitTitle:   pulumi.String("PR_TITLE"),
		Template:                 t,
		Topics:                   pulumi.ToStringArray(repo.Topics),
		Visibility:               pulumi.String(repo.Visibility),
		VulnerabilityAlerts:      inputs.Yes,
		WebCommitSignoffRequired: inputs.No,
	}

	return github.NewRepository(ctx, repo.Name, args)
}

func template(repo *Repository) (*github.RepositoryTemplateArgs, error) {
	if !repo.HasTemplate() {
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
	if !repo.HasPages() {
		return nil
	}

	return &github.RepositoryPagesArgs{
		BuildType: pulumi.String("legacy"),
		Cname:     pulumi.String(repo.Pages.CNAME),
		Source: &github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(master),
		},
	}
}
