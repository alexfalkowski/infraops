package app

import (
	"fmt"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func httpProbe(path string) cv1.ProbeArgs {
	return cv1.ProbeArgs{
		HttpGet: cv1.HTTPGetActionArgs{
			Path: pulumi.String(path),
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds: pulumi.Int(5),
		PeriodSeconds:       pulumi.Int(10),
		SuccessThreshold:    pulumi.Int(1),
		FailureThreshold:    pulumi.Int(3),
		TimeoutSeconds:      pulumi.Int(30),
	}
}

func tcpProbe() cv1.ProbeArgs {
	return cv1.ProbeArgs{
		TcpSocket: cv1.TCPSocketActionArgs{
			Port: pulumi.Int(8080),
		},
		InitialDelaySeconds:           pulumi.Int(5),
		PeriodSeconds:                 pulumi.Int(10),
		SuccessThreshold:              pulumi.Int(1),
		FailureThreshold:              pulumi.Int(5),
		TimeoutSeconds:                pulumi.Int(60),
		TerminationGracePeriodSeconds: pulumi.Int(30),
	}
}

func probePath(name, path string) string {
	return fmt.Sprintf("/%s/%s", name, path)
}
