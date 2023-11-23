package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createTemplate(ctx *pulumi.Context) error {
	_, err := github.NewRepository(ctx, "go-service-template", &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		AllowUpdateBranch:   pulumi.Bool(true),
		DefaultBranch:       pulumi.String("master"),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String("A template for go services"),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		IsTemplate:          pulumi.Bool(true),
		Name:                pulumi.String("go-service-template"),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("enabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("enabled"),
			},
		},
		SquashMergeCommitTitle: pulumi.String("PR_TITLE"),
		Visibility:             pulumi.String("public"),
		VulnerabilityAlerts:    pulumi.Bool(true),
	})

	return err
}
