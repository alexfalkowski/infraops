package gh

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateRepositories provisions repositories and their settings, ordering managed templates
// before their consumers regardless of configuration order. Templates owned by alexfalkowski
// and included in repos use resource outputs; other templates must already exist.
// Empty input creates no resources. Duplicate names and template cycles return errors.
func CreateRepositories(ctx *pulumi.Context, repos []*Repository) error {
	creator := &repositoryCreator{
		managed:  make(map[string]*Repository, len(repos)),
		created:  make(map[string]*github.Repository, len(repos)),
		creating: make(map[string]bool, len(repos)),
	}
	for _, repo := range repos {
		name := strings.ToLower(repo.Name)
		if _, ok := creator.managed[name]; ok {
			return fmt.Errorf("duplicate repository %s", repo.Name)
		}
		creator.managed[name] = repo
	}

	for _, repo := range repos {
		if err := creator.create(ctx, repo); err != nil {
			return err
		}
	}
	return nil
}

type repositoryCreator struct {
	managed  map[string]*Repository
	created  map[string]*github.Repository
	creating map[string]bool
}

func (c *repositoryCreator) create(ctx *pulumi.Context, repo *Repository) error {
	name := strings.ToLower(repo.Name)
	if _, ok := c.created[name]; ok {
		return nil
	}
	if c.creating[name] {
		return fmt.Errorf("repository %s: template dependency cycle", repo.Name)
	}
	c.creating[name] = true

	var source *github.Repository
	if repo.HasTemplate() && strings.EqualFold(repo.Template.Owner, "alexfalkowski") {
		templateName := strings.ToLower(repo.Template.Repository)
		if template, ok := c.managed[templateName]; ok {
			if err := c.create(ctx, template); err != nil {
				return err
			}
			source = c.created[templateName]
		}
	}

	r, err := createRepository(ctx, source, repo)
	if err != nil {
		return err
	}
	c.created[name] = r
	delete(c.creating, name)
	return nil
}
