package version

import (
	"errors"
	"slices"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/app"
)

// ErrNotFound is returned by Update when the named application does not exist in the configuration.
var ErrNotFound = errors.New("name not found")

// Update sets the version for the application identified by name in the configuration at path.
//
// The configuration is loaded using app.ReadConfiguration and persisted using app.WriteConfiguration.
// If no application matches name, Update returns ErrNotFound.
func Update(name, version, path string) error {
	config, err := app.ReadConfiguration(path)
	if err != nil {
		return err
	}

	i := slices.IndexFunc(config.GetApplications(), func(a *v2.Application) bool {
		return a.GetName() == name
	})
	if i == -1 {
		return ErrNotFound
	}

	config.GetApplications()[i].Version = version
	return app.WriteConfiguration(path, config)
}
