package app

import (
	"errors"
	"strings"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErrVersionMismatch indicates that an application's deployed version does not match the
// configured version.
//
// Note: this error is defined for callers to use as a sentinel; it is not necessarily
// produced by functions in this file.
var ErrVersionMismatch = errors.New("version mismatch")

// ReadConfiguration reads the Kubernetes applications configuration from path.
//
// The file is expected to be HJSON matching the v2.Kubernetes protobuf schema.
func ReadConfiguration(path string) (*v2.Kubernetes, error) {
	var configuration v2.Kubernetes
	err := config.Read(path, &configuration)
	return &configuration, err
}

// WriteConfiguration writes configuration to path in the repository's canonical HJSON form.
//
// The underlying writer preserves the destination file mode and appends a trailing newline.
func WriteConfiguration(path string, configuration *v2.Kubernetes) error {
	return config.Write(path, configuration)
}

// ConvertApplication converts a protobuf v2.Application into the internal App model.
//
// This conversion normalizes optional fields (for example env vars) into Go-friendly
// structures that are later used to create Kubernetes resources.
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

// CreateApplication creates all Kubernetes resources required for app in the target cluster.
//
// Resources are created in a fixed order to ensure dependencies exist (for example,
// ServiceAccount before Deployment). Any error aborts the creation flow and is returned.
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

// App describes an application that will be deployed to Kubernetes.
//
// It is derived from v2.Application and is used as the input to resource constructors
// within this package.
type App struct {
	// Resources are optional resource requests/limits to apply to the pod.
	Resources *Resources
	// ID is an optional identifier from configuration.
	ID string
	// Name is the Kubernetes application name (used for resource naming).
	Name string
	// Kind determines how the application is deployed (for example "internal" vs "external").
	Kind string
	// Namespace is the Kubernetes namespace to deploy into.
	Namespace string
	// Domain is the external hostname associated with the application.
	Domain string
	// Version is the application version string used for deployment (for example as an image tag).
	Version string
	// Secrets is a list of secret names referenced by this application.
	Secrets []string
	// EnvVars are environment variables to be injected into the application container.
	EnvVars []*EnvVar
}

// HasResources reports whether a has resource requirements configured.
func (a *App) HasResources() bool {
	return a.Resources != nil
}

// IsInternal reports whether this application uses the repository's opinionated deployment model.
//
// Internal applications are typically built/published by this repository and deployed
// using a conventional container image naming scheme.
func (a *App) IsInternal() bool {
	return a.Kind == "internal"
}

// IsExternal reports whether this application is not built/published by this repository.
func (a *App) IsExternal() bool {
	return a.Kind == "external"
}

// Resources describes optional CPU/memory/storage ranges for an application's pod.
type Resources struct {
	CPU     *Range
	Memory  *Range
	Storage *Range
}

// Range represents a min/max range expressed as Kubernetes quantity strings.
type Range struct {
	Min string
	Max string
}

// EnvVar represents an environment variable to inject into the application container.
type EnvVar struct {
	Name  string
	Value string
}

// IsSecret reports whether the env var value refers to a secret.
//
// Secret values are encoded using the "secret:" prefix and are resolved by this package
// into Kubernetes Secret references during resource creation.
func (e *EnvVar) IsSecret() bool {
	return strings.HasPrefix(e.Value, "secret:")
}
