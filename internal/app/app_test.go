package app_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/internal/app"
	test "github.com/alexfalkowski/infraops/internal/pulumi"
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

	_, err := app.Read("invalid")
	require.Error(t, err)
}

func withResource(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "1234",
		Name:          "test",
		Namespace:     "test",
		Domain:        "test.com",
		InitVersion:   "1.0.0",
		Version:       "1.0.0",
		ConfigVersion: "1.0.0",
		Secrets:       []string{"test"},
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
	}

	return app.Create(ctx, a)
}

func withoutResource(ctx *pulumi.Context) error {
	a := &app.App{
		ID:            "1234",
		Name:          "test",
		Namespace:     "test",
		Domain:        "test.com",
		InitVersion:   "1.0.0",
		Version:       "1.0.0",
		ConfigVersion: "1.0.0",
	}

	return app.Create(ctx, a)
}
