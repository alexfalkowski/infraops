package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func collaborator(ctx *pulumi.Context, repo *Repository) error {
	if !repo.EnableCollaborators {
		return nil
	}

	args := &github.RepositoryCollaboratorArgs{
		Permission: pulumi.String("admin"),
		Repository: pulumi.String("alexfalkowski/" + repo.Name),
		Username:   pulumi.String("lean-thoughts-ci"),
	}

	_, err := github.NewRepositoryCollaborator(ctx, repo.Name, args)

	return err
}
