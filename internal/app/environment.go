package app

import (
	"strings"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func addEnvironments(app *App, envs cv1.EnvVarArray) cv1.EnvVarArray {
	for _, envVar := range app.EnvVars {
		var arg cv1.EnvVarArgs

		if envVar.IsSecret() {
			value := strings.TrimPrefix(envVar.Value, "secret:")
			name, value, _ := strings.Cut(value, "/")

			arg = cv1.EnvVarArgs{
				Name: pulumi.String(envVar.Name),
				ValueFrom: &cv1.EnvVarSourceArgs{
					SecretKeyRef: &cv1.SecretKeySelectorArgs{
						Name: pulumi.String(name + secretSuffix),
						Key:  pulumi.String(value),
					},
				},
			}
		} else {
			arg = cv1.EnvVarArgs{
				Name:  pulumi.String(envVar.Name),
				Value: pulumi.String(envVar.Value),
			}
		}

		envs = append(envs, arg)
	}

	return envs
}
