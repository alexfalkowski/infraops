package app

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErrVersionMismatch for app.
var ErrVersionMismatch = errors.New("version mismatch")

type (
	// App to be created.
	App struct {
		Resources     *Resources
		ID            string
		Name          string
		Namespace     string
		Domain        string
		InitVersion   string
		Version       string
		ConfigVersion string
		Secrets       Secrets
	}

	// Secrets for apps.
	Secrets []string

	// Resources for apps.
	Resources struct {
		CPU     *Range
		Memory  *Range
		Storage *Range
	}

	// Range for apps.
	Range struct {
		Min string
		Max string
	}
)

// CreateApp in the cluster.
func CreateApp(ctx *pulumi.Context, app *App) error {
	fns := []func(ctx *pulumi.Context, app *App) error{
		createServiceAccount, createNetworkPolicy,
		createConfigMap, createPodDisruptionBudget,
		createDeployment, createService, createIngress,
	}

	for _, fn := range fns {
		if err := fn(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

// HasConfigVersion for app.
func (a *App) HasConfigVersion() bool {
	return a.ConfigVersion != ""
}

// HasResources for app.
func (a *App) HasResources() bool {
	return a.Resources != nil
}
