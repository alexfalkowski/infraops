package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/gh"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := gh.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			repositories := config.GetRepositories()
			run := func(ctx *pulumi.Context) error {
				for _, repository := range repositories {
					if err := gh.CreateRepository(ctx, gh.ConvertRepository(repository)); err != nil {
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
		{name: "area", path: "gh.hjson"},
		{name: "shared", path: filepath.Join("..", "..", "internal", "test", "gh.hjson")},
	}
}

func runWithMocks(run pulumi.RunFunc, mocks pulumi.MockResourceMonitor) error {
	return pulumi.RunErr(run, pulumi.WithMocks("project", "stack", mocks))
}
