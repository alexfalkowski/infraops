package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func podSecurity() cv1.PodSecurityContextArgs {
	return cv1.PodSecurityContextArgs{
		FsGroup:      pulumi.Int(2000),
		RunAsNonRoot: pulumi.Bool(true),
		RunAsUser:    pulumi.Int(10001),
		RunAsGroup:   pulumi.Int(10001),
	}
}
