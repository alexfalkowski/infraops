package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Service for apps.
type Service struct {
	Name      string
	ClusterIP string
}

func createService(ctx *pulumi.Context, app *App) error {
	if app.IsExternal() {
		args := &cv1.ServiceArgs{
			Metadata: metadata(app.Name, app, matchLabels(app)),
			Spec: cv1.ServiceSpecArgs{
				Ports: cv1.ServicePortArray{
					servicePort("http", 8080),
				},
				Selector: matchLabels(app),
				Type:     pulumi.String("ClusterIP"),
			},
		}

		_, err := cv1.NewService(ctx, app.Name, args)

		return err
	}

	services := []*Service{
		{Name: app.Name},
		{Name: "headless-" + app.Name, ClusterIP: "None"},
	}

	for _, service := range services {
		ports := cv1.ServicePortArray{
			servicePort("debug", 6060),
			servicePort("http", 8080),
			servicePort("grpc", 9090),
		}
		ip := service.ClusterIP
		spec := cv1.ServiceSpecArgs{
			ClusterIP: pulumi.String(ip),
			Ports:     ports,
			Type:      pulumi.String("ClusterIP"),
		}

		if len(ip) > 0 {
			spec.Selector = matchLabels(app)
		}

		args := &cv1.ServiceArgs{
			Metadata: metadata(service.Name, app, matchLabels(app)),
			Spec:     spec,
		}

		_, err := cv1.NewService(ctx, service.Name, args)
		if err != nil {
			return err
		}
	}

	return nil
}
