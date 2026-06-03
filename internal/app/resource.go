package app

import (
	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createResources(app *App) cv1.ResourceRequirementsArgs {
	if !app.HasResources() {
		return cv1.ResourceRequirementsArgs{}
	}

	requests := pulumi.StringMap{}
	limits := pulumi.StringMap{}

	if app.Resources.CPU != nil {
		requests["cpu"] = pulumi.String(app.Resources.CPU.Min)
		limits["cpu"] = pulumi.String(app.Resources.CPU.Max)
	}
	if app.Resources.Memory != nil {
		requests["memory"] = pulumi.String(app.Resources.Memory.Min)
		limits["memory"] = pulumi.String(app.Resources.Memory.Max)
	}
	if app.Resources.Storage != nil {
		requests["ephemeral-storage"] = pulumi.String(app.Resources.Storage.Min)
		limits["ephemeral-storage"] = pulumi.String(app.Resources.Storage.Max)
	}
	return cv1.ResourceRequirementsArgs{
		Requests: requests,
		Limits:   limits,
	}
}

var resources = ResourcesMap{
	"small": {
		CPU:     &Range{Min: "125m", Max: "250m"},
		Memory:  &Range{Min: "64Mi", Max: "128Mi"},
		Storage: &Range{Min: "1Gi", Max: "2Gi"},
	},
}

// ResourcesMap maps configuration resource profile names to Kubernetes request/limit ranges.
type ResourcesMap map[string]*Resources

// Resources returns the named resource profile, or the small profile when name is unknown.
func (r ResourcesMap) Resources(name string) *Resources {
	res, ok := r[name]
	if ok {
		return res
	}
	return r["small"]
}
