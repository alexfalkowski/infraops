package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := app.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			applications := config.GetApplications()
			run := func(ctx *pulumi.Context) error {
				for _, application := range applications {
					if err := app.CreateApplication(ctx, app.ConvertApplication(application)); err != nil {
						return err
					}
				}

				return nil
			}

			require.NoError(t, runWithMocks(run, &test.Stub{}))
			require.Error(t, runWithMocks(run, &test.ErrStub{}))
		})
	}
}

type fixture struct {
	name string
	path string
}

func fixtures() []fixture {
	return []fixture{
		{name: "area", path: "apps.hjson"},
		{name: "shared", path: filepath.Join("..", "..", "internal", "test", "apps.hjson")},
	}
}

func runWithMocks(run pulumi.RunFunc, mocks pulumi.MockResourceMonitor) error {
	return pulumi.RunErr(run, pulumi.WithMocks("project", "stack", mocks))
}
