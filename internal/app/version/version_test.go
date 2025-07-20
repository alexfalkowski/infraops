package version_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app/version"
	"github.com/stretchr/testify/require"
)

func TestValidUpdate(t *testing.T) {
	err := version.Update("app1", "1.1.0", "../test/apps.pbtxt")
	require.NoError(t, err)
}

func TestInvalidUpdate(t *testing.T) {
	err := version.Update("app1", "1.1.0", "../test/none.pbtxt")
	require.Error(t, err)
}

func TestMissingUpdate(t *testing.T) {
	err := version.Update("none", "1.1.0", "../test/apps.pbtxt")
	require.Error(t, err)
}
