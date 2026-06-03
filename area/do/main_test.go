package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestCreateCluster(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			stub.FailAllResources()

			require.NoError(t, test.RunWithMocks(run(fixture.path), &test.Stub{}))
			require.Error(t, test.RunWithMocks(run(fixture.path), stub))
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
