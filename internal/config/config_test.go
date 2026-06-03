package config_test

import (
	"path/filepath"
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestWriteReturnsStatError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.hjson")

	require.Error(t, config.Write(path, &v2.Kubernetes{}))
}
