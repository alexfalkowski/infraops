package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestCreateRepository(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			stub.FailAllResources()

			require.NoError(t, test.RunWithMocks(run(fixture.path), &test.Stub{}))
			require.Error(t, test.RunWithMocks(run(fixture.path), stub))
		})
	}
}

func TestCreateRepositoryDoesNotCreateFromFailedTemplate(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailAllResources()

	require.Error(t, test.RunWithMocks(run("test/templates.hjson"), stub))
	repositories := stub.Resources("github:index/repository:Repository")
	require.Len(t, repositories, 1, "the consumer must wait for its template to succeed")
	require.Equal(t, "starter", test.Property(t, repositories[0], "name").StringValue())
}

type fixture struct {
	name string
	path string
}

func fixtures() []fixture {
	return []fixture{
		{name: "area", path: "gh.hjson"},
		{name: "shared", path: filepath.Join("..", "..", "internal", "test", "gh.hjson")},
		{name: "managed template", path: "test/templates.hjson"},
	}
}
