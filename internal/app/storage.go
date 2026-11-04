package app

import (
	"errors"
	"fmt"
	"path"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const storageSuffix = "-data"

// Storage configures one persistent filesystem claim for an application.
// Nil App.Storage preserves the existing config/secret-only volumes.
type Storage struct {
	// Profile preserves nil (cluster default) versus empty (no class).
	Profile *string

	// Capacity is a required Kubernetes storage request; the provider validates quantities.
	Capacity string

	// Path is a required absolute, clean, non-root directory mounted writable.
	// The volume must support the pod's existing fsGroup 2000 permissions.
	Path string
}

// Validate checks the required capacity and absolute, clean, non-root data path.
// Quantity validity and profile availability are checked by Kubernetes/provider operations.
func (s *Storage) Validate() error {
	if s.Capacity == "" {
		return errors.New("storage capacity is required")
	}
	if !path.IsAbs(s.Path) || s.Path == "/" || path.Clean(s.Path) != s.Path {
		return errors.New("storage path must be an absolute, clean, non-root directory")
	}
	return nil
}

func createStorage(ctx *pulumi.Context, app *App) error {
	if app.Storage == nil {
		return nil
	}
	storage := app.Storage
	if err := storage.Validate(); err != nil {
		return fmt.Errorf("application %s: %w", app.Name, err)
	}

	m := metadata(app)
	m.Name = pulumi.String(app.Name + storageSuffix)
	spec := cv1.PersistentVolumeClaimSpecArgs{
		AccessModes: pulumi.StringArray{pulumi.String("ReadWriteOnce")},
		VolumeMode:  pulumi.String("Filesystem"),
		Resources: cv1.VolumeResourceRequirementsArgs{
			Requests: pulumi.StringMap{"storage": pulumi.String(storage.Capacity)},
		},
	}
	if storage.Profile != nil {
		spec.StorageClassName = pulumi.String(*storage.Profile)
	}

	args := &cv1.PersistentVolumeClaimArgs{
		Metadata: m,
		Spec:     spec,
	}

	// Replacing or deleting a data claim needs deliberate backup/migration and unprotection.
	_, err := cv1.NewPersistentVolumeClaim(ctx, app.Name+storageSuffix, args, pulumi.Protect(true))
	return err
}
