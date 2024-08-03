package app_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	fns := []pulumi.RunFunc{withResource, withoutResource}
	for _, f := range fns {
		require.NoError(t, pulumi.RunErr(f, pulumi.WithMocks("project", "stack", app.Mocks(0))))
	}
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
		SecretVolumes: []string{"test"},
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
	}

	return app.CreateApp(ctx, a)
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

	return app.CreateApp(ctx, a)
}
