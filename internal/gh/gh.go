package gh

import (
	"errors"
	"fmt"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const master = "master"

var (
	// ErrMissingTemplate for gh.
	ErrMissingTemplate = errors.New("missing template")

	// ErrMissingChecks for gh.
	ErrMissingChecks = errors.New("missing checks")

	// Public visibility.
	Public = Visibility("public")

	// Private visibility.
	Private = Visibility("private")
)

// Template for gh.
type Template struct {
	Owner      string
	Repository string
}

// Valid if no error is returned.
func (t *Template) Valid() error {
	if t.Owner == "" || t.Repository == "" {
		return ErrMissingTemplate
	}
	return nil
}

// Checks for gh.
type Checks []string

// Valid if no error is returned.
func (c Checks) Valid() error {
	if len(c) == 0 {
		return ErrMissingChecks
	}
	return nil
}

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Github, error) {
	var configuration v2.Github
	err := config.Read(path, &configuration)
	return &configuration, err
}

// ConvertRepository converts a v2.Repository to a Repository.
func ConvertRepository(r *v2.Repository) *Repository {
	repository := &Repository{
		Name:        r.GetName(),
		Description: r.GetDescription(),
		HomepageURL: r.GetHomepageUrl(),
		Visibility:  Visibility(r.GetVisibility()),
		Topics:      r.GetTopics(),
		Checks:      Checks(r.GetChecks()),
		IsTemplate:  r.GetIsTemplate(),
		Archived:    r.GetArchived(),
	}

	if collaborators := r.GetCollaborators(); collaborators != nil {
		repository.Collaborators = &Collaborators{
			Enabled: collaborators.GetEnabled(),
		}
	}

	if template := r.GetTemplate(); template != nil {
		owner := template.GetOwner()
		repo := template.GetRepository()

		if owner != "" && repo != "" {
			repository.Template = &Template{
				Owner:      owner,
				Repository: repo,
			}
		}
	}

	if pages := r.GetPages(); pages != nil {
		repository.Pages = &Pages{
			Enabled: pages.GetEnabled(),
			CNAME:   pages.GetCname(),
		}
	}

	return repository
}

// CreateRepository for gh.
func CreateRepository(ctx *pulumi.Context, repo *Repository) error {
	r, err := repository(ctx, repo)
	if err != nil {
		return fmt.Errorf("%v: %w", repo.Name, err)
	}

	if err := branchProtection(ctx, r.NodeId, repo); err != nil {
		return fmt.Errorf("%v: %w", repo.Name, err)
	}

	if err := collaborator(ctx, repo); err != nil {
		return fmt.Errorf("%v: %w", repo.Name, err)
	}

	return nil
}

type (
	// Visibility of the repositories.
	Visibility string

	// Collaborators for gh.
	Collaborators struct {
		Enabled bool
	}

	// Pages for gh.
	Pages struct {
		CNAME   string
		Enabled bool
	}
)

// Repository for gh.
type Repository struct {
	Collaborators *Collaborators
	Template      *Template
	Pages         *Pages
	Name          string
	Description   string
	HomepageURL   string
	Visibility    Visibility
	Topics        []string
	Checks        Checks
	IsTemplate    bool
	Archived      bool
}

// HasCollaborators for this repository.
func (r *Repository) HasCollaborators() bool {
	return r.Collaborators != nil && r.Collaborators.Enabled
}

// HasTemplate for this repository.
func (r *Repository) HasTemplate() bool {
	return r.Template != nil
}

// HasPages for this repository.
func (r *Repository) HasPages() bool {
	return r.Pages != nil && r.Pages.Enabled
}
