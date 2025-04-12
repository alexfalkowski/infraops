package app

import (
	"strings"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func addEnvironments(app *App, envs cv1.EnvVarArray) cv1.EnvVarArray {
	for _, env := range app.Environments {
		var arg cv1.EnvVarArgs

		if env.IsSecret() {
			value := strings.TrimPrefix(env.Value, "secret:")
			name, value, _ := strings.Cut(value, "/")

			arg = cv1.EnvVarArgs{
				Name: pulumi.String(env.Name),
				ValueFrom: &cv1.EnvVarSourceArgs{
					SecretKeyRef: &cv1.SecretKeySelectorArgs{
						Name: pulumi.String(name + secretSuffix),
						Key:  pulumi.String(value),
					},
				},
			}
		} else {
			arg = cv1.EnvVarArgs{
				Name:  pulumi.String(env.Name),
				Value: pulumi.String(env.Value),
			}
		}

		envs = append(envs, arg)
	}

	return envs
}
