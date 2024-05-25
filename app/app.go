package app

import (
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// KonfigVersion used by apps.
const KonfigVersion = "v1.131.3"

// App to be created.
type App struct {
	Name          string
	Version       string
	ConfigVersion string
	SecretVolumes []string
}

// CreateApp in the cluster.
func CreateApp(ctx *pulumi.Context, app *App) error {
	err := createServiceAccount(ctx, app)
	if err != nil {
		return err
	}

	err = createNetworkPolicy(ctx, app)
	if err != nil {
		return err
	}

	err = createConfigMap(ctx, app)
	if err != nil {
		return err
	}

	err = createDeployment(ctx, app)
	if err != nil {
		return err
	}

	err = createService(ctx, app)
	if err != nil {
		return err
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
