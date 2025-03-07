package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/internal/app"
	test "github.com/alexfalkowski/infraops/internal/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateApp(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, application := range applications {
			err := app.CreateApp(ctx, application)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, application := range applications {
			err := app.CreateApp(ctx, application)
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
