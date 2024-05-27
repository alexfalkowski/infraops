package app

import (
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type (
	// App to be created.
	App struct {
		Name          string
		Version       string
		ConfigVersion string
		InitVersion   string
		Memory        Memory
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

func metadata(app *App) mv1.ObjectMetaArgs {
	return mv1.ObjectMetaArgs{
		Name:      pulumi.String(app.Name),
		Namespace: pulumi.String(app.Name),
	}
}

func labels(app *App) pulumi.StringMap {
	return pulumi.StringMap{"app": pulumi.String(app.Name)}
}
