package version_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/app/version"
	"github.com/stretchr/testify/require"
)

func TestValidUpdate(t *testing.T) {
	path := copyFixture(t)

	require.NoError(t, version.Update("app1", "2.0.0", path))

	config, err := app.ReadConfiguration(path)
	require.NoError(t, err)

	applications := config.GetApplications()
	require.Len(t, applications, 2)
	require.Equal(t, "app1", applications[0].GetName())
	require.Equal(t, "2.0.0", applications[0].GetVersion())
	require.EqualValues(t, 3, applications[0].GetReplicas())
	require.Equal(t, "app2", applications[1].GetName())
	require.Equal(t, "1.0.0", applications[1].GetVersion())
	require.EqualValues(t, 3, applications[1].GetReplicas())
}

func TestInvalidUpdate(t *testing.T) {
	err := version.Update("app1", "1.1.0", "../test/none.hjson")
	require.Error(t, err)
}

func TestMissingUpdate(t *testing.T) {
	err := version.Update("none", "1.1.0", "../test/apps.hjson")
	require.ErrorIs(t, err, version.ErrNotFound)
}

func copyFixture(t *testing.T) string {
	t.Helper()

	data, err := os.ReadFile("../test/apps.hjson")
	require.NoError(t, err)

	file, err := os.CreateTemp(t.TempDir(), "apps-*.hjson")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, file.Close())
	}()

	_, err = file.Write(data)
	require.NoError(t, err)

	return filepath.Clean(file.Name())
}
