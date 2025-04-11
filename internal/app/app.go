package app

import (
	"errors"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErrVersionMismatch for app.
var ErrVersionMismatch = errors.New("version mismatch")

type (
	// App to be created.
	App struct {
		Resources     *Resources
		ID            string
		Name          string
		Kind          string
		Namespace     string
		Domain        string
		InitVersion   string
		Version       string
		ConfigVersion string
		Secrets       Secrets
	}

	// Secrets for apps.
	Secrets []string

	// Resources for apps.
	Resources struct {
		CPU     *Range
		Memory  *Range
		Storage *Range
	}

	// Range for apps.
	Range struct {
		Min string
		Max string
	}
)

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Kubernetes, error) {
	var configuration v2.Kubernetes
	err := config.Read(path, &configuration)

	return &configuration, err
}

// ConvertApplication converts a v2.Application to an App.
func ConvertApplication(a *v2.Application) *App {
	app := &App{
		ID:            a.GetId(),
		Kind:          a.GetKind(),
		Name:          a.GetName(),
		Namespace:     a.GetNamespace(),
		Domain:        a.GetDomain(),
		InitVersion:   a.GetInitVersion(),
		Version:       a.GetVersion(),
		ConfigVersion: a.GetConfigVersion(),
		Secrets:       a.GetSecrets(),
	}

	resources := a.GetResources()
	if resources != nil {
		cpu := resources.GetCpu()
		mem := resources.GetMemory()
		storage := resources.GetStorage()
		r := &Resources{}

		if cpu != nil {
			r.CPU = &Range{Min: cpu.GetMin(), Max: cpu.GetMax()}
		}

		if mem != nil {
			r.Memory = &Range{Min: mem.GetMin(), Max: mem.GetMax()}
		}

		if storage != nil {
			r.Storage = &Range{Min: storage.GetMin(), Max: storage.GetMax()}
		}

		app.Resources = r
	}

	return app
}

// CreateApplication in the cluster.
func CreateApplication(ctx *pulumi.Context, app *App) error {
	fns := []func(ctx *pulumi.Context, app *App) error{
		createServiceAccount, createNetworkPolicy,
		createConfigMap, createPodDisruptionBudget,
		createDeployment, createService, createIngress,
	}

	for _, fn := range fns {
		if err := fn(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

// HasConfigVersion for app.
func (a *App) HasConfigVersion() bool {
	return a.ConfigVersion != ""
}

// HasResources for app.
func (a *App) HasResources() bool {
	return a.Resources != nil
}

// IsInternal defines whether this application uses our opinionated way of deploying apps.
func (a *App) IsInternal() bool {
	return a.Kind == "internal"
}
