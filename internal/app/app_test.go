package app_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	tests := []struct {
		app  *app.App
		name string
	}{
		{name: "with resource", app: withResource()},
		{name: "without resource", app: withoutResource()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &resourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
			require.Equal(t, portNumbers(app.Ports(tt.app)), networkPolicyIngressPorts(stub.networkPolicy))
		})
	}

	for _, tt := range tests {
		t.Run(tt.name+" error", func(t *testing.T) {
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.Error(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", &test.ErrStub{})))
		})
	}

	_, err := app.ReadConfiguration("invalid")
	require.Error(t, err)
}

type resourceStub struct {
	test.Stub

	networkPolicy resource.PropertyMap
}

func (s *resourceStub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	if args.TypeToken == "kubernetes:networking.k8s.io/v1:NetworkPolicy" {
		s.networkPolicy = args.Inputs
	}

	return s.Stub.NewResource(args)
}

func withResource() *app.App {
	return &app.App{
		ID:        "1234",
		Kind:      "internal",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Secrets:   []string{"test"},
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		EnvVars: []*app.EnvVar{
			{Name: "test", Value: "test"},
		},
	}
}

func withoutResource() *app.App {
	return &app.App{
		ID:        "1234",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Secrets:   []string{"test"},
		EnvVars: []*app.EnvVar{
			{Name: "test", Value: "test"},
		},
	}
}

func portNumbers(ports []app.Port) []int {
	numbers := make([]int, 0, len(ports))
	for _, port := range ports {
		numbers = append(numbers, port.Number)
	}

	return numbers
}

func networkPolicyIngressPorts(policy resource.PropertyMap) []int {
	spec := policy[resource.PropertyKey("spec")].ObjectValue()
	ingress := spec[resource.PropertyKey("ingress")].ArrayValue()
	ports := ingress[0].ObjectValue()[resource.PropertyKey("ports")].ArrayValue()

	values := make([]int, 0, len(ports))
	for _, port := range ports {
		value := port.ObjectValue()[resource.PropertyKey("port")].NumberValue()
		values = append(values, int(value))
	}

	return values
}
