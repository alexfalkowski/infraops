package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createService(ctx *pulumi.Context, app *App) error {
	args := &cv1.ServiceArgs{
		Metadata: metadata(app, matchLabels(app)),
		Spec: cv1.ServiceSpecArgs{
			Ports:    servicePorts(app),
			Selector: matchLabels(app),
			Type:     pulumi.String("ClusterIP"),
		},
	}
	_, err := cv1.NewService(ctx, app.Name, args)

	return err
}
