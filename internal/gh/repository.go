package gh

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// repository creates or updates a GitHub repository and applies the baseline repository settings.
//
// It validates repo.Checks (required for branch protection) and, when configured, attaches a
// repository template and GitHub Pages configuration.
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
		AllowSquashMerge:    inputs.Yes,
		AllowUpdateBranch:   inputs.Yes,
		AutoInit:            inputs.Yes,
		AllowAutoMerge:      inputs.Yes,
		DeleteBranchOnMerge: inputs.Yes,
		Description:         pulumi.String(repo.Description),
		HasIssues:           inputs.Yes,
		HasProjects:         inputs.Yes,
		HasWiki:             inputs.Yes,
		HomepageUrl:         pulumi.String(repo.HomepageURL),
		IsTemplate:          pulumi.Bool(repo.IsTemplate),
		Name:                pulumi.String(repo.Name),
		Archived:            pulumi.Bool(repo.Archived),
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
		WebCommitSignoffRequired: inputs.No,
	}

	repository, err := github.NewRepository(ctx, repo.Name, args)
	if err != nil {
		return nil, err
	}

	if err := vulnerabilityAlerts(ctx, repo); err != nil {
		return nil, err
	}

	if err := pages(ctx, repo); err != nil {
		return nil, err
	}

	return repository, nil
}

func vulnerabilityAlerts(ctx *pulumi.Context, repo *Repository) error {
	_, err := github.NewRepositoryVulnerabilityAlerts(ctx, repo.Name, &github.RepositoryVulnerabilityAlertsArgs{
		Enabled:    inputs.Yes,
		Repository: pulumi.String(repo.Name),
	})
	return err
}

func pages(ctx *pulumi.Context, repo *Repository) error {
	if !repo.HasPages() {
		return nil
	}

	_, err := github.NewRepositoryPages(ctx, repo.Name, &github.RepositoryPagesArgs{
		BuildType:  pulumi.String("legacy"),
		Cname:      pulumi.String(repo.Pages.CNAME),
		Repository: pulumi.String(repo.Name),
		Source: &github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(master),
		},
	})
	return err
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
