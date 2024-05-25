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
			PodSelector: mv1.LabelSelectorArgs{MatchLabels: labels(app)},
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
