package app

import (
	"errors"
	"strings"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErrVersionMismatch for app.
var ErrVersionMismatch = errors.New("version mismatch")

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Kubernetes, error) {
	var configuration v2.Kubernetes
	err := config.Read(path, &configuration)
	return &configuration, err
}

// WriteConfiguration writes the configuration to a file.
func WriteConfiguration(path string, configuration *v2.Kubernetes) error {
	return config.Write(path, configuration)
}

// ConvertApplication converts a v2.Application to an App.
func ConvertApplication(a *v2.Application) *App {
	app := &App{
		ID: a.GetId(), Kind: a.GetKind(),
		Name: a.GetName(), Namespace: a.GetNamespace(),
		Domain: a.GetDomain(), Version: a.GetVersion(),
		Resources: resources.Resources(a.GetResource()),
		Secrets:   a.GetSecrets(),
	}

	envVars := a.GetEnvVars()
	if envVars != nil {
		app.EnvVars = make([]*EnvVar, len(envVars))

		for i, e := range envVars {
			envVar := &EnvVar{
				Name:  e.GetName(),
				Value: e.GetValue(),
			}

			app.EnvVars[i] = envVar
		}
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

// App to be created.
type App struct {
	Resources *Resources
	ID        string
	Name      string
	Kind      string
	Namespace string
	Domain    string
	Version   string
	Secrets   []string
	EnvVars   []*EnvVar
}

// HasResources for app.
func (a *App) HasResources() bool {
	return a.Resources != nil
}

// IsInternal defines whether this application uses our opinionated way of deploying apps.
func (a *App) IsInternal() bool {
	return a.Kind == "internal"
}

// IsExternal defines an app that is not built by us.
func (a *App) IsExternal() bool {
	return a.Kind == "external"
}

// Resources for apps.
type Resources struct {
	CPU     *Range
	Memory  *Range
	Storage *Range
}

// Range for apps.
type Range struct {
	Min string
	Max string
}

// EnvVar for apps.
type EnvVar struct {
	Name  string
	Value string
}

// IsSecret defines whether the env variable is a secret.
func (e *EnvVar) IsSecret() bool {
	return strings.HasPrefix(e.Value, "secret:")
}
