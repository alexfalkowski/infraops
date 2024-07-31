package app_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/app"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		a := &app.App{
			ID:            "1234",
			Name:          "test",
			Namespace:     "test",
			Domain:        "test.com",
			InitVersion:   "1.0.0",
			Version:       "1.0.0",
			ConfigVersion: "1.0.0",
			SecretVolumes: []string{"test"},
			Memory:        app.Memory{Min: "64Mi", Max: "128Mi"},
		}

		err := app.CreateApp(ctx, a)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", app.Mocks(0)))

	require.NoError(t, err)
}
