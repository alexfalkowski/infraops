package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// collaborator provisions a GitHub repository collaborator when collaborator management is enabled.
//
// The current implementation grants "admin" permissions to the "lean-thoughts-ci" user for the
// repository identified by repo.Name under the "alexfalkowski" owner.
func collaborator(ctx *pulumi.Context, repo *Repository) error {
	if !repo.HasCollaborators() {
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
