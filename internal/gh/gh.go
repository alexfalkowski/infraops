package gh

import (
	"errors"
	"fmt"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// master is the default branch name used by this repository's GitHub configuration.
const master = "master"

var (
	// ErrMissingTemplate is returned when a template repository configuration is incomplete.
	ErrMissingTemplate = errors.New("missing template")

	// ErrMissingChecks is returned when a repository is configured to require status checks
	// but no checks are listed.
	ErrMissingChecks = errors.New("missing checks")

	// Public is GitHub "public" repository visibility.
	Public = Visibility("public")

	// Private is GitHub "private" repository visibility.
	Private = Visibility("private")
)

// Template identifies a GitHub template repository used to create a new repository.
type Template struct {
	// Owner is the GitHub organization or user that owns the template repository.
	Owner string
	// Repository is the template repository name.
	Repository string
}

// Valid validates that the template repository reference is complete.
//
// If either Owner or Repository is empty, Valid returns ErrMissingTemplate.
func (t *Template) Valid() error {
	if t.Owner == "" || t.Repository == "" {
		return ErrMissingTemplate
	}
	return nil
}

// Checks is a list of required GitHub status check contexts used by branch protection.
type Checks []string

// Valid validates that there is at least one check context.
//
// If the list is empty, Valid returns ErrMissingChecks.
func (c Checks) Valid() error {
	if len(c) == 0 {
		return ErrMissingChecks
	}
	return nil
}

// ReadConfiguration reads the GitHub area configuration from path.
//
// The file is expected to be HJSON matching the v2.Github protobuf schema.
func ReadConfiguration(path string) (*v2.Github, error) {
	var configuration v2.Github
	err := config.Read(path, &configuration)
	return &configuration, err
}

// ConvertRepository converts a protobuf v2.Repository into the internal Repository model.
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

// CreateRepository provisions and configures a GitHub repository.
//
// It creates (or updates) the repository and then applies additional repository configuration
// such as branch protection and collaborators. Errors are wrapped with the repository name
// to make failures easier to identify in Pulumi output.
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
	// Visibility is a GitHub repository visibility setting (for example "public" or "private").
	Visibility string

	// Collaborators describes whether collaborator management is enabled for a repository.
	Collaborators struct {
		// Enabled controls whether collaborators should be managed for the repository.
		Enabled bool
	}

	// Pages describes GitHub Pages configuration for a repository.
	Pages struct {
		// CNAME is the custom domain to configure for GitHub Pages (optional).
		CNAME string
		// Enabled controls whether Pages should be managed/enabled for the repository.
		Enabled bool
	}
)

// Repository describes a GitHub repository and its desired configuration.
type Repository struct {
	// Collaborators is optional; when nil or disabled, collaborator resources are not managed.
	Collaborators *Collaborators
	// Template is optional; when set, it identifies the template repository used on creation.
	Template *Template
	// Pages is optional; when nil or disabled, Pages resources are not managed.
	Pages *Pages

	// Name is the repository name.
	Name string
	// Description is the repository description.
	Description string
	// HomepageURL is the repository homepage URL.
	HomepageURL string
	// Visibility controls repository visibility.
	Visibility Visibility
	// Topics are repository topics to apply.
	Topics []string
	// Checks are required status checks to enforce via branch protection.
	Checks Checks
	// IsTemplate marks this repository as a template repository.
	IsTemplate bool
	// Archived controls whether the repository should be archived.
	Archived bool
}

// HasCollaborators reports whether collaborator management is enabled for this repository.
func (r *Repository) HasCollaborators() bool {
	return r.Collaborators != nil && r.Collaborators.Enabled
}

// HasTemplate reports whether this repository is configured to be created from a template.
func (r *Repository) HasTemplate() bool {
	return r.Template != nil
}

// HasPages reports whether Pages management is enabled for this repository.
func (r *Repository) HasPages() bool {
	return r.Pages != nil && r.Pages.Enabled
}
