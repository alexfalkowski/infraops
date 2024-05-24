package app

import (
	"os"
	"path/filepath"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createConfigMap(ctx *pulumi.Context, app *App) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	p := filepath.Clean(filepath.Join(wd, app.Name+".yaml"))

	d, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	args := &cv1.ConfigMapArgs{
		Metadata: metadata(app),
		Data:     pulumi.StringMap{app.Config: pulumi.String(string(d))},
	}
	_, err = cv1.NewConfigMap(ctx, app.Name, args)

	return err
}
