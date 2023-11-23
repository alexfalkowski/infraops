package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createAppConfig(ctx *pulumi.Context) error {
	_, err := github.NewRepository(ctx, "app-config", &github.RepositoryArgs{
		AllowUpdateBranch:   pulumi.Bool(true),
		DefaultBranch:       pulumi.String("master"),
		DeleteBranchOnMerge: pulumi.Bool(true),
		Description:         pulumi.String("A place for all of my application configuration"),
		HasDownloads:        pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		Name:                pulumi.String("app-config"),
		SecurityAndAnalysis: &github.RepositorySecurityAndAnalysisArgs{
			SecretScanning: &github.RepositorySecurityAndAnalysisSecretScanningArgs{
				Status: pulumi.String("disabled"),
			},
			SecretScanningPushProtection: &github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
				Status: pulumi.String("disabled"),
			},
		},
		Visibility: pulumi.String("public"),
	})

	return err
}
