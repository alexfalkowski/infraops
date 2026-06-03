package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			stub.FailAllResources()

			require.NoError(t, test.RunWithMocks(run(fixture.path), &test.Stub{}))
			require.Error(t, test.RunWithMocks(run(fixture.path), stub))
		})
	}
}

func TestLeanMakefileCoversConfiguredApps(t *testing.T) {
	config, err := app.ReadConfiguration("apps.hjson")
	require.NoError(t, err)

	content, err := os.ReadFile("lean.mk")
	require.NoError(t, err)

	targets := makefileTargets(string(content))
	for _, workflow := range []string{"rollout", "verify", "load"} {
		aggregate := workflow + "-lean"
		require.Contains(t, targets, aggregate)

		for _, application := range config.GetApplications() {
			target := workflow + "-" + application.GetName()
			require.Contains(t, targets[aggregate], target)
			require.Contains(t, targets, target)
		}
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

func makefileTargets(content string) map[string][]string {
	targets := map[string][]string{}
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "@") {
			continue
		}

		target, dependencies, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		targets[strings.TrimSpace(target)] = strings.Fields(dependencies)
	}

	return targets
}
