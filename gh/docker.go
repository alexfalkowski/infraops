package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDocker(ctx *pulumi.Context) error {
	_, err := github.NewRepository(ctx, "docker", &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		DefaultBranch:       pulumi.String("master"),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String("Common setup used for my projects."),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		Name:                pulumi.String("docker"),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("enabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("enabled"),
			},
		},
		Topics: pulumi.StringArray{
			pulumi.String("hbase"),
			pulumi.String("docker"),
			pulumi.String("kubernetes"),
			pulumi.String("ruby"),
			pulumi.String("golang"),
		},
		Visibility:          pulumi.String("public"),
		VulnerabilityAlerts: pulumi.Bool(true),
	})

	return err
}
