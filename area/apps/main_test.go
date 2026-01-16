package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	config, err := app.ReadConfiguration("apps.hjson")
	require.NoError(t, err)

	applications := config.GetApplications()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, application := range applications {
			err := app.CreateApplication(ctx, app.ConvertApplication(application))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, application := range applications {
			err := app.CreateApplication(ctx, app.ConvertApplication(application))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
