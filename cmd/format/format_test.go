package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatRejectsMismatchedConfigKind(t *testing.T) {
	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})

	tests := []struct {
		name     string
		kind     string
		contents []byte
	}{
		{
			name:     "apps config as cloudflare",
			kind:     "cf",
			contents: []byte("{\n  version: \"2.0\"\n  applications: []\n}\n"),
		},
		{
			name:     "cloudflare config as digital ocean",
			kind:     "do",
			contents: []byte("{\n  version: \"2.0\"\n  balancer_zones: []\n}\n"),
		},
		{
			name:     "digital ocean config as github",
			kind:     "gh",
			contents: []byte("{\n  version: \"2.0\"\n  clusters: []\n}\n"),
		},
		{
			name:     "github config as apps",
			kind:     "apps",
			contents: []byte("{\n  version: \"2.0\"\n  repositories: []\n}\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.hjson")
			require.NoError(t, os.WriteFile(path, test.contents, 0o600))

			os.Args = []string{"format", "-k", test.kind, "-p", path}

			require.Error(t, run())

			actual, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, test.contents, actual)
		})
	}
}
