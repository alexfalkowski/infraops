package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateCluster(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := do.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			clusters := config.GetClusters()
			run := func(ctx *pulumi.Context) error {
				for _, cluster := range clusters {
					if err := do.CreateCluster(ctx, do.ConvertCluster(cluster)); err != nil {
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
		{name: "area", path: "do.hjson"},
		{name: "shared", path: filepath.Join("..", "..", "internal", "test", "do.hjson")},
	}
}

func runWithMocks(run pulumi.RunFunc, mocks pulumi.MockResourceMonitor) error {
	return pulumi.RunErr(run, pulumi.WithMocks("project", "stack", mocks))
}
