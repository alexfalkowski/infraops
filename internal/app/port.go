package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func servicePort(name string, port int) cv1.ServicePortArgs {
	return cv1.ServicePortArgs{
		AppProtocol: pulumi.String("TCP"),
		Name:        pulumi.String(name),
		Port:        pulumi.Int(port),
		TargetPort:  pulumi.Int(port),
	}
}
