package app

import (
	"errors"
	"fmt"
	"path"
	"regexp"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const storageClass = "do-block-storage-retain"

var (
	storageKey  = regexp.MustCompile(`^[a-z0-9]+$`)
	storageName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// Storage configures a persistent filesystem claim for an application.
type Storage struct {
	// Key is the configured volume key, appended to the application name to name the claim.
	Key string

	// Capacity is a required Kubernetes storage request; the provider validates quantities.
	Capacity string

	// Path is a required absolute, clean, non-root directory mounted writable.
	// The volume must support the pod's existing fsGroup 2000 permissions.
	Path string
}

func (s Storage) claimName(appName string) string {
	return appName + "-" + s.Key
}

func (s Storage) validate(name string) error {
	// Alphanumeric keys make the final hyphen the separator between app and key.
	if !storageKey.MatchString(s.Key) {
		return errors.New("storage key must use only lowercase letters and digits")
	}
	if len(name) > 63 || !storageName.MatchString(name) {
		return errors.New("storage name must be a Kubernetes DNS label of at most 63 characters")
	}
	if s.Capacity == "" {
		return errors.New("storage capacity is required")
	}
	if !path.IsAbs(s.Path) || s.Path == "/" || path.Clean(s.Path) != s.Path {
		return errors.New("storage path must be an absolute, clean, non-root directory")
	}
	return nil
}

func createStorage(ctx *pulumi.Context, app *App) ([]pulumi.Resource, error) {
	if err := validateStorage(app); err != nil {
		return nil, err
	}
	claims := make([]pulumi.Resource, 0, len(app.Storage))
	for _, storage := range app.Storage {
		claim, err := createStorageClaim(ctx, app, storage)
		if err != nil {
			return nil, err
		}
		claims = append(claims, claim)
	}
	return claims, nil
}

func validateStorage(app *App) error {
	if len(app.Storage) == 0 {
		return nil
	}
	names := make(map[string]bool, len(app.Secrets)+len(app.Storage))
	for _, secret := range app.Secrets {
		names[secret+secretSuffix] = true
	}
	paths := make(map[string]bool, len(app.Storage))
	for _, storage := range app.Storage {
		name := storage.claimName(app.Name)
		if err := storage.validate(name); err != nil {
			return fmt.Errorf("application %s, volume %q: %w", app.Name, storage.Key, err)
		}
		if names[name] {
			return fmt.Errorf("application %s: duplicate volume name %q", app.Name, name)
		}
		if paths[storage.Path] {
			return fmt.Errorf("application %s: duplicate storage path %q", app.Name, storage.Path)
		}
		names[name] = true
		paths[storage.Path] = true
	}
	return nil
}

func createStorageClaim(ctx *pulumi.Context, app *App, storage Storage) (*cv1.PersistentVolumeClaim, error) {
	name := storage.claimName(app.Name)
	m := metadata(app)
	m.Name = pulumi.String(name)
	spec := cv1.PersistentVolumeClaimSpecArgs{
		AccessModes:      pulumi.StringArray{pulumi.String("ReadWriteOnce")},
		VolumeMode:       pulumi.String("Filesystem"),
		StorageClassName: pulumi.String(storageClass),
		Resources: cv1.VolumeResourceRequirementsArgs{
			Requests: pulumi.StringMap{"storage": pulumi.String(storage.Capacity)},
		},
	}

	args := &cv1.PersistentVolumeClaimArgs{
		Metadata: m,
		Spec:     spec,
	}

	// Replacing or deleting a data claim needs deliberate backup/migration and unprotection.
	return cv1.NewPersistentVolumeClaim(ctx, name, args, pulumi.Protect(true))
}
