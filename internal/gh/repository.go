package gh

import (
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var repositoryLabels = []struct {
	name        string
	color       string
	description string
}{
	{
		name:        "automerge",
		color:       "0E8A16",
		description: "Eligible for automated merging.",
	},
	{
		name:        "hold",
		color:       "FBCA04",
		description: "Do not merge until this label is removed.",
	},
}

// repository creates or updates a GitHub repository and applies the baseline repository settings.
//
// It validates repo.Checks (required for branch protection) and, when configured, attaches a
// repository template and GitHub Pages configuration. Secret scanning, push protection, and
// vulnerability alerts are enabled for every repository. Non-archived repositories receive
// the shared repository labels.
func repository(ctx *pulumi.Context, source *github.Repository, repo *Repository) (*github.Repository, error) {
	templateArgs, err := template(source, repo)
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
		Template:                 templateArgs,
		Topics:                   pulumi.ToStringArray(repo.Topics),
		Visibility:               pulumi.String(repo.Visibility),
		WebCommitSignoffRequired: inputs.No,
	}

	repository, err := github.NewRepository(ctx, repo.Name, args)
	if err != nil {
		return nil, err
	}

	// Use repository outputs so dependent settings wait for repository creation.
	if err := labels(ctx, repository, repo); err != nil {
		return nil, err
	}

	if err := vulnerabilityAlerts(ctx, repository, repo); err != nil {
		return nil, err
	}

	if err := pages(ctx, repository, repo); err != nil {
		return nil, err
	}

	return repository, nil
}

func labels(ctx *pulumi.Context, created *github.Repository, repo *Repository) error {
	if repo.Archived {
		return nil
	}

	for _, label := range repositoryLabels {
		_, err := github.NewIssueLabel(ctx, repo.Name+"-"+label.name, &github.IssueLabelArgs{
			Repository:  created.Name,
			Name:        pulumi.String(label.name),
			Color:       pulumi.String(label.color),
			Description: pulumi.String(label.description),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func vulnerabilityAlerts(ctx *pulumi.Context, created *github.Repository, repo *Repository) error {
	_, err := github.NewRepositoryVulnerabilityAlerts(ctx, repo.Name, &github.RepositoryVulnerabilityAlertsArgs{
		Enabled:    inputs.Yes,
		Repository: created.Name,
	})
	return err
}

func pages(ctx *pulumi.Context, created *github.Repository, repo *Repository) error {
	if !repo.HasPages() {
		return nil
	}

	_, err := github.NewRepositoryPages(ctx, repo.Name, &github.RepositoryPagesArgs{
		BuildType:  pulumi.String("legacy"),
		Cname:      pulumi.String(repo.Pages.CNAME),
		Repository: created.Name,
		Source: &github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(master),
		},
	})
	return err
}

func template(source *github.Repository, repo *Repository) (*github.RepositoryTemplateArgs, error) {
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
	if source != nil {
		args.Repository = source.Name
	}

	return args, nil
}
