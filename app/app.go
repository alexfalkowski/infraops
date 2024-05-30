package app

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type (
	// App to be created.
	App struct {
		Memory        Memory
		ID            string
		Name          string
		Domain        string
		Version       string
		ConfigVersion string
		InitVersion   string
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
