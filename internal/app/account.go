package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createServiceAccount(ctx *pulumi.Context, app *App) error {
	args := &cv1.ServiceAccountArgs{Metadata: metadata(app)}
	_, err := cv1.NewServiceAccount(ctx, app.Name, args)

	return err
}
