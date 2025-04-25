package gh

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const bot = "lean-thoughts-ci"

func collaborator(ctx *pulumi.Context, repo *Repository) error {
	args := &github.RepositoryCollaboratorArgs{
		Permission: pulumi.String("admin"),
		Repository: pulumi.String("alexfalkowski/" + repo.Name),
		Username:   pulumi.String(bot),
	}

	_, err := github.NewRepositoryCollaborator(ctx, repo.Name, args)

	return err
}
