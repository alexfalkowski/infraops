package version

import (
	"errors"
	"slices"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/app"
)

// ErrNotFound for version.
var ErrNotFound = errors.New("name not found")

// Update name with a version by providing the path to the config.
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
