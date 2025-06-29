package app_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	fns := []pulumi.RunFunc{withResource, withoutResource}

	for _, f := range fns {
		require.NoError(t, pulumi.RunErr(f, pulumi.WithMocks("project", "stack", &test.Stub{})))
	}

	for _, f := range fns {
		require.Error(t, pulumi.RunErr(f, pulumi.WithMocks("project", "stack", &test.ErrStub{})))
	}

	_, err := app.ReadConfiguration("invalid")
	require.Error(t, err)
}

func withResource(ctx *pulumi.Context) error {
	a := &app.App{
		ID:        "1234",
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

	return app.CreateApplication(ctx, a)
}

func withoutResource(ctx *pulumi.Context) error {
	a := &app.App{
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

	return app.CreateApplication(ctx, a)
}
