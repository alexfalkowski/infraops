package app

import (
	nv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createIngress(ctx *pulumi.Context, app *App) error {
	if app.ConfigVersion == "" {
		return nil
	}

	args := &nv1.IngressArgs{
		Metadata: metadata(app),
		Spec: nv1.IngressSpecArgs{
			IngressClassName: pulumi.String("nginx"),
			Rules: nv1.IngressRuleArray{
				nv1.IngressRuleArgs{
					Host: pulumi.String(app.Name + ".lean-thoughts.com"),
					Http: httpIngressRule(app),
				},
			},
		},
	}
	_, err := nv1.NewIngress(ctx, app.Name, args)

	return err
}

func httpIngressRule(app *App) nv1.HTTPIngressRuleValueArgs {
	return nv1.HTTPIngressRuleValueArgs{
		Paths: nv1.HTTPIngressPathArray{
			nv1.HTTPIngressPathArgs{
				Path:     pulumi.String("/"),
				PathType: pulumi.String("Prefix"),
				Backend: nv1.IngressBackendArgs{
					Service: nv1.IngressServiceBackendArgs{
						Name: pulumi.String(app.Name),
						Port: nv1.ServiceBackendPortArgs{Number: pulumi.Int(8080)},
					},
				},
			},
		},
	}
}
