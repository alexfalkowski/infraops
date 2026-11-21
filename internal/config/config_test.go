package config_test

import (
	"path/filepath"
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/stretchr/testify/require"
)

func TestReadRejectsDuplicateVolumeKeys(t *testing.T) {
	err := config.Read("test/duplicate-volume-key.hjson", &v2.Kubernetes{})
	require.ErrorContains(t, err, "duplicate HJSON object key")
}

func TestReadRejectsMixedCaseStorageFields(t *testing.T) {
	err := config.Read("test/duplicate-storage-field.hjson", &v2.Kubernetes{})
	require.ErrorContains(t, err, `duplicate HJSON field "storage"`)
}

func TestReadAcceptsSingleMixedCaseStorageField(t *testing.T) {
	var cfg v2.Kubernetes
	require.NoError(t, config.Read("test/storage-field-alias.hjson", &cfg))
	require.Len(t, cfg.GetApplications(), 1)
	require.Equal(t, "/data", cfg.GetApplications()[0].GetStorage()["data"].GetPath())
}

func TestReadOmitsValuesFromDuplicateKeyError(t *testing.T) {
	err := config.Read("test/duplicate-env-value.hjson", &v2.Kubernetes{})
	require.ErrorContains(t, err, "duplicate HJSON object key")
	require.NotContains(t, err.Error(), "first-synthetic-value")
	require.NotContains(t, err.Error(), "second-synthetic-value")
	require.NotContains(t, err.Error(), ">>>")
}

func TestWriteReturnsStatError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.hjson")

	require.Error(t, config.Write(path, &v2.Kubernetes{}))
}
