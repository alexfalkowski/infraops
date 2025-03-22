package app

import (
	"errors"

	v1 "github.com/alexfalkowski/infraops/api/infraops/v1"
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
func ReadConfiguration(path string) (*v1.Kubernetes, error) {
	var configuration v1.Kubernetes
	err := config.Read(path, &configuration)

	return &configuration, err
}

// ConvertApplication converts a v1.Application to an App.
func ConvertApplication(a *v1.Application) *App {
	cpu := a.GetResources().GetCpu()
	mem := a.GetResources().GetMemory()
	storage := a.GetResources().GetStorage()
	app := &App{
		ID:            a.GetId(),
		Name:          a.GetName(),
		Namespace:     a.GetNamespace(),
		Domain:        a.GetDomain(),
		InitVersion:   a.GetInitVersion(),
		Version:       a.GetVersion(),
		ConfigVersion: a.GetConfigVersion(),
		Secrets:       a.GetSecrets(),
		Resources: &Resources{
			CPU:     &Range{Min: cpu.GetMin(), Max: cpu.GetMax()},
			Memory:  &Range{Min: mem.GetMin(), Max: mem.GetMax()},
			Storage: &Range{Min: storage.GetMin(), Max: storage.GetMax()},
		},
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
