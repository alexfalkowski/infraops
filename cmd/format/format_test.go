package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatRejectsMismatchedConfigKind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cf.hjson")
	contents := []byte("{\n  version: \"2.0\"\n  balancer_zones: []\n}\n")
	require.NoError(t, os.WriteFile(path, contents, 0o600))

	args := os.Args
	t.Cleanup(func() {
		os.Args = args
	})
	os.Args = []string{"format", "-k", "apps", "-p", path}

	require.Error(t, run())

	actual, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, contents, actual)
}
