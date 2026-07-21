package app_test

import (
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/app"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

const (
	configMapResourceType      = "kubernetes:core/v1:ConfigMap"
	deploymentResourceType     = "kubernetes:apps/v1:Deployment"
	ingressResourceType        = "kubernetes:networking.k8s.io/v1:Ingress"
	networkPolicyResourceType  = "kubernetes:networking.k8s.io/v1:NetworkPolicy"
	pdbResourceType            = "kubernetes:policy/v1:PodDisruptionBudget"
	serviceAccountResourceType = "kubernetes:core/v1:ServiceAccount"
	serviceResourceType        = "kubernetes:core/v1:Service"
)

func TestApp(t *testing.T) {
	tests := []struct {
		app  *app.App
		name string
	}{
		{name: "with resource", app: appWithResources()},
		{name: "without resource", app: appWithoutResources()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
			require.Equal(t, portNumbers(app.Ports(tt.app)), networkPolicyIngressPorts(t, resourceOf(t, stub, networkPolicyResourceType)))
			require.Equal(t, portProtocols(app.Ports(tt.app)), servicePortAppProtocols(t, resourceOf(t, stub, serviceResourceType)))
		})
	}

	for _, tt := range tests {
		t.Run(tt.name+" error", func(t *testing.T) {
			stub := &test.ResourceStub{}
			stub.FailResource(serviceAccountResourceType)

			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.Error(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
		})
	}

	_, err := app.ReadConfiguration("invalid")
	require.Error(t, err)
}

func TestApplicationReturnsResourceErrors(t *testing.T) {
	for _, tt := range []struct {
		name         string
		resourceType string
	}{
		{name: "service account", resourceType: serviceAccountResourceType},
		{name: "network policy", resourceType: networkPolicyResourceType},
		{name: "config map", resourceType: configMapResourceType},
		{name: "pod disruption budget", resourceType: pdbResourceType},
		{name: "deployment", resourceType: deploymentResourceType},
		{name: "service", resourceType: serviceResourceType},
		{name: "ingress", resourceType: ingressResourceType},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			stub.FailResource(tt.resourceType)

			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, appWithResources())
			}

			require.Error(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
			require.NotEmpty(t, stub.Resources(tt.resourceType))
		})
	}
}

func TestApplicationAllowsExplicitZeroReplicas(t *testing.T) {
	stub := &test.ResourceStub{}
	application := app.ConvertApplication(&v2.Application{
		Kind:      "external",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Replicas:  0,
	})

	run := func(ctx *pulumi.Context) error {
		return app.CreateApplication(ctx, application)
	}

	require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	deployment := resourceOf(t, stub, deploymentResourceType)
	require.Zero(t, test.Property(t, deploymentSpec(deployment), "replicas").NumberValue())
}

func TestApplicationReturnsMissingConfigError(t *testing.T) {
	stub := &test.ResourceStub{}
	application := appWithResources()
	application.Namespace = "missing-config-fixture"

	run := func(ctx *pulumi.Context) error {
		return app.CreateApplication(ctx, application)
	}

	require.Error(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	require.Len(t, stub.Resources(serviceAccountResourceType), 1)
	require.Len(t, stub.Resources(networkPolicyResourceType), 1)
	require.Empty(t, stub.Resources(configMapResourceType))
}

func TestConvertedApplicationDeploymentInputs(t *testing.T) {
	stub := &test.ResourceStub{}
	application := app.ConvertApplication(&v2.Application{
		Id:        "1234",
		Kind:      "internal",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Resource:  "unknown",
		Replicas:  4,
		Secrets:   []string{"database"},
		EnvVars: []*v2.EnvVar{
			{Name: "LOG_LEVEL", Value: "info"},
			{Name: "DATABASE_PASSWORD", Value: "secret:database/password"},
		},
	})
	run := func(ctx *pulumi.Context) error {
		return app.CreateApplication(ctx, application)
	}

	require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	deployment := resourceOf(t, stub, deploymentResourceType)
	require.Equal(t, "125m", resourceRequest(deployment, "cpu"))
	require.Equal(t, "250m", resourceLimit(deployment, "cpu"))
	require.Equal(t, "64Mi", resourceRequest(deployment, "memory"))
	require.Equal(t, "128Mi", resourceLimit(deployment, "memory"))
	require.Equal(t, "1Gi", resourceRequest(deployment, "ephemeral-storage"))
	require.Equal(t, "2Gi", resourceLimit(deployment, "ephemeral-storage"))
	require.Equal(t, "RuntimeDefault", seccompProfileType(deployment))
	require.Equal(t, 4, int(test.Property(t, deploymentSpec(deployment), "replicas").NumberValue()))

	secret := envVar(deployment, "DATABASE_PASSWORD")
	valueFrom := test.Property(t, secret, "valueFrom").ObjectValue()
	secretKeyRef := test.Property(t, valueFrom, "secretKeyRef").ObjectValue()
	require.Equal(t, "database-secret", test.Property(t, secretKeyRef, "name").StringValue())
	require.Equal(t, "password", test.Property(t, secretKeyRef, "key").StringValue())
}

func TestApplicationContainerSecurityContextDropsAllCapabilities(t *testing.T) {
	tests := []struct {
		app  *app.App
		name string
	}{
		{app: appWithResources(), name: "internal"},
		{
			app: app.ConvertApplication(&v2.Application{
				Kind:      "external",
				Name:      "test",
				Namespace: "test",
				Domain:    "test.com",
				Version:   "1.0.0",
				Replicas:  1,
			}),
			name: "external",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
			require.Equal(t, []string{"ALL"}, droppedCapabilities(t, resourceOf(t, stub, deploymentResourceType)))
		})
	}
}

func TestApplicationDelaysShutdownForRoutingPropagation(t *testing.T) {
	tests := []struct {
		app  *app.App
		name string
	}{
		{app: appWithResources(), name: "internal"},
		{
			app: app.ConvertApplication(&v2.Application{
				Kind:      "external",
				Name:      "test",
				Namespace: "test",
				Domain:    "test.com",
				Version:   "1.0.0",
				Replicas:  1,
			}),
			name: "external",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, tt.app)
			}

			require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
			deployment := resourceOf(t, stub, deploymentResourceType)
			require.Equal(t, 35, int(test.Property(t, podSpec(deployment), "terminationGracePeriodSeconds").NumberValue()))
			require.Equal(t, 5, preStopSleepSeconds(t, deployment))
		})
	}
}

func TestExternalApplicationOmitsInternalResources(t *testing.T) {
	stub := &test.ResourceStub{}
	application := app.ConvertApplication(&v2.Application{
		Id:        "1234",
		Kind:      "external",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Resource:  "small",
		Replicas:  1,
	})
	run := func(ctx *pulumi.Context) error {
		return app.CreateApplication(ctx, application)
	}

	require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	deployment := resourceOf(t, stub, deploymentResourceType)
	componentLabel := resource.PropertyKey("circleci.com/component-name")
	require.Empty(t, stub.Resources(configMapResourceType))
	require.Empty(t, test.Property(t, metadata(deployment), "annotations").ObjectValue())
	require.NotContains(t, test.Property(t, metadata(deployment), "labels").ObjectValue(), componentLabel)
	require.NotContains(t, test.Property(t, podTemplateMetadata(deployment), "labels").ObjectValue(), componentLabel)
	require.NotContains(t, podSpec(deployment), resource.PropertyKey("volumes"))
	require.NotContains(t, container(deployment), resource.PropertyKey("volumeMounts"))
}

func appWithResources() *app.App {
	return &app.App{
		ID:        "1234",
		Kind:      "internal",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Replicas:  3,
		Secrets:   []string{"test"},
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		EnvVars: []*app.EnvVar{
			{Name: "test", Value: "test"},
		},
	}
}

func appWithoutResources() *app.App {
	return &app.App{
		ID:        "1234",
		Name:      "test",
		Namespace: "test",
		Domain:    "test.com",
		Version:   "1.0.0",
		Replicas:  3,
		Secrets:   []string{"test"},
		EnvVars: []*app.EnvVar{
			{Name: "test", Value: "test"},
		},
	}
}

func portNumbers(ports []app.Port) []int {
	numbers := make([]int, 0, len(ports))
	for _, port := range ports {
		numbers = append(numbers, port.Number)
	}

	return numbers
}

func portProtocols(ports []app.Port) []string {
	protocols := make([]string, 0, len(ports))
	for _, port := range ports {
		protocols = append(protocols, port.Protocol)
	}

	return protocols
}

func resourceOf(t *testing.T, stub *test.ResourceStub, token string) resource.PropertyMap {
	t.Helper()

	resources := stub.Resources(token)
	require.Len(t, resources, 1)

	return resources[0]
}

func networkPolicyIngressPorts(t *testing.T, policy resource.PropertyMap) []int {
	t.Helper()

	spec := test.Property(t, policy, "spec").ObjectValue()
	ingress := test.Property(t, spec, "ingress").ArrayValue()
	ports := test.Property(t, ingress[0].ObjectValue(), "ports").ArrayValue()

	values := make([]int, 0, len(ports))
	for _, port := range ports {
		value := port.ObjectValue()[resource.PropertyKey("port")].NumberValue()
		values = append(values, int(value))
	}

	return values
}

func metadata(deployment resource.PropertyMap) resource.PropertyMap {
	return deployment[resource.PropertyKey("metadata")].ObjectValue()
}

func podTemplateMetadata(deployment resource.PropertyMap) resource.PropertyMap {
	return deploymentSpec(deployment)[resource.PropertyKey("template")].ObjectValue()[resource.PropertyKey("metadata")].ObjectValue()
}

func podSpec(deployment resource.PropertyMap) resource.PropertyMap {
	template := deploymentSpec(deployment)[resource.PropertyKey("template")].ObjectValue()
	return template[resource.PropertyKey("spec")].ObjectValue()
}

func container(deployment resource.PropertyMap) resource.PropertyMap {
	containers := podSpec(deployment)[resource.PropertyKey("containers")].ArrayValue()
	return containers[0].ObjectValue()
}

func envVar(deployment resource.PropertyMap, name string) resource.PropertyMap {
	envs := container(deployment)[resource.PropertyKey("env")].ArrayValue()
	for _, env := range envs {
		value := env.ObjectValue()
		if value[resource.PropertyKey("name")].StringValue() == name {
			return value
		}
	}

	return nil
}

func resourceRequest(deployment resource.PropertyMap, name string) string {
	return resourceValue(deployment, "requests", name)
}

func resourceLimit(deployment resource.PropertyMap, name string) string {
	return resourceValue(deployment, "limits", name)
}

func resourceValue(deployment resource.PropertyMap, kind, name string) string {
	resources := container(deployment)[resource.PropertyKey("resources")].ObjectValue()
	values := resources[resource.PropertyKey(kind)].ObjectValue()
	return values[resource.PropertyKey(name)].StringValue()
}

func seccompProfileType(deployment resource.PropertyMap) string {
	security := container(deployment)[resource.PropertyKey("securityContext")].ObjectValue()
	profile := security[resource.PropertyKey("seccompProfile")].ObjectValue()
	return profile[resource.PropertyKey("type")].StringValue()
}

func droppedCapabilities(t *testing.T, deployment resource.PropertyMap) []string {
	t.Helper()

	security := container(deployment)[resource.PropertyKey("securityContext")].ObjectValue()
	capabilities := test.Property(t, security, "capabilities").ObjectValue()
	return test.StringValues(test.Property(t, capabilities, "drop").ArrayValue())
}

func preStopSleepSeconds(t *testing.T, deployment resource.PropertyMap) int {
	t.Helper()

	lifecycle := test.Property(t, container(deployment), "lifecycle").ObjectValue()
	preStop := test.Property(t, lifecycle, "preStop").ObjectValue()
	require.NotContains(t, preStop, resource.PropertyKey("exec"))
	sleep := test.Property(t, preStop, "sleep").ObjectValue()

	return int(test.Property(t, sleep, "seconds").NumberValue())
}

func deploymentSpec(deployment resource.PropertyMap) resource.PropertyMap {
	return deployment[resource.PropertyKey("spec")].ObjectValue()
}

func servicePortAppProtocols(t *testing.T, service resource.PropertyMap) []string {
	t.Helper()

	spec := test.Property(t, service, "spec").ObjectValue()
	ports := test.Property(t, spec, "ports").ArrayValue()

	values := make([]string, 0, len(ports))
	for _, port := range ports {
		value := port.ObjectValue()[resource.PropertyKey("appProtocol")].StringValue()
		values = append(values, value)
	}

	return values
}
