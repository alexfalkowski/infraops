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
		Memory        Memory
		ID            string
		Name          string
		Namespace     string
		Domain        string
		InitVersion   string
		Version       string
		ConfigVersion string
		SecretVolumes []string
	}

	// Memory for apps.
	Memory struct {
		Min string
		Max string
	}

	createFn func(ctx *pulumi.Context, app *App) error
)

// CreateApp in the cluster.
func CreateApp(ctx *pulumi.Context, app *App) error {
	fns := []createFn{
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
