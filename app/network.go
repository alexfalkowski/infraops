package app

import (
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	nv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createNetworkPolicy(ctx *pulumi.Context, app *App) error {
	args := &nv1.NetworkPolicyArgs{
		Metadata: metadata(app),
		Spec: nv1.NetworkPolicySpecArgs{
			PodSelector: mv1.LabelSelectorArgs{MatchLabels: matchLabels(app)},
			Ingress: nv1.NetworkPolicyIngressRuleArray{
				nv1.NetworkPolicyIngressRuleArgs{},
			},
			Egress: nv1.NetworkPolicyEgressRuleArray{
				nv1.NetworkPolicyEgressRuleArgs{},
			},
		},
	}
	_, err := nv1.NewNetworkPolicy(ctx, app.Name, args)

	return err
}

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
					Host: pulumi.String(app.Name + "." + app.Domain),
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
