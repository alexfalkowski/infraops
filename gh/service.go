//nolint:dupl
package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createService(ctx *pulumi.Context) error {
	_, err := github.NewRepository(ctx, "go-service", &github.RepositoryArgs{
		AllowMergeCommit:    pulumi.Bool(false),
		AllowRebaseMerge:    pulumi.Bool(false),
		AllowUpdateBranch:   pulumi.Bool(true),
		DefaultBranch:       pulumi.String("master"),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String("A framework to build services in go."),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		Name:                pulumi.String("go-service"),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("enabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("enabled"),
			},
		},
		Visibility:          pulumi.String("public"),
		VulnerabilityAlerts: pulumi.Bool(true),
	})

	return err
}
