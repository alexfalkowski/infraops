package app_test

import (
	"os"
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
	pvcResourceType            = "kubernetes:core/v1:PersistentVolumeClaim"
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
			ingress := resourceOf(t, stub, ingressResourceType)
			spec := test.Property(t, ingress, "spec").ObjectValue()
			require.Equal(t, "contour", test.Property(t, spec, "ingressClassName").StringValue())
			tls := test.Property(t, spec, "tls").ArrayValue()
			require.Len(t, tls, 1)
			require.Equal(t, "ingress-tls", test.Property(t, tls[0].ObjectValue(), "secretName").StringValue())
			hosts := test.StringValues(test.Property(t, tls[0].ObjectValue(), "hosts").ArrayValue())
			require.Equal(t, []string{tt.app.Name + "." + tt.app.Domain}, hosts)
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
		Kind:      "worker",
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
	application.Storage = &app.Storage{Capacity: "1Gi", Path: "/data"}

	run := func(ctx *pulumi.Context) error {
		return app.CreateApplication(ctx, application)
	}

	require.Error(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	require.Len(t, stub.Resources(serviceAccountResourceType), 1)
	require.Len(t, stub.Resources(networkPolicyResourceType), 1)
	require.Empty(t, stub.Resources(configMapResourceType))
	require.Empty(t, stub.Resources(pvcResourceType), "missing application config must not allocate persistent storage")
}

func TestConvertedApplicationDeploymentInputs(t *testing.T) {
	stub := &test.ResourceStub{}
	application := app.ConvertApplication(&v2.Application{
		Id:        "1234",
		Kind:      "api",
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
		{app: appWithResources(), name: "api"},
		{
			app: app.ConvertApplication(&v2.Application{
				Kind:      "worker",
				Name:      "test",
				Namespace: "test",
				Domain:    "test.com",
				Version:   "1.0.0",
				Replicas:  1,
			}),
			name: "worker",
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
		{app: appWithResources(), name: "api"},
		{
			app: app.ConvertApplication(&v2.Application{
				Kind:      "worker",
				Name:      "test",
				Namespace: "test",
				Domain:    "test.com",
				Version:   "1.0.0",
				Replicas:  1,
			}),
			name: "worker",
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

func TestWorkerApplicationOmitsServiceAndIngress(t *testing.T) {
	stub := &test.ResourceStub{}
	application := app.ConvertApplication(&v2.Application{
		Id:        "1234",
		Kind:      "worker",
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
	require.NotEmpty(t, stub.Resources(deploymentResourceType))
	require.Empty(t, stub.Resources(serviceResourceType))
	require.Empty(t, stub.Resources(ingressResourceType))
}

func TestApplicationPersistentStorage(t *testing.T) {
	for _, tt := range []struct {
		storage *v2.Storage
		name    string
	}{
		{name: "without storage"},
		{name: "cluster default class", storage: &v2.Storage{Capacity: "1Gi", Path: "/data"}},
		{
			name: "no class",
			storage: &v2.Storage{
				Capacity: "2Gi", Path: "/var/lib/test", Profile: new(""),
			},
		},
		{
			name: "named class",
			storage: &v2.Storage{
				Capacity: "1Gi", Path: "/data", Profile: new("example-csi"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			configuration, err := app.ReadConfiguration("test/storage.hjson")
			require.NoError(t, err)
			configured := configuration.GetApplications()[0]
			configured.Storage = tt.storage
			// Canonical config rewrites must retain the distinction between absent and empty class.
			file := t.TempDir() + "/apps.hjson"
			require.NoError(t, os.WriteFile(file, nil, 0o600))
			require.NoError(t, app.WriteConfiguration(file, configuration))
			configuration, err = app.ReadConfiguration(file)
			require.NoError(t, err)
			configured = configuration.GetApplications()[0]
			application := app.ConvertApplication(configured)
			stub := &test.ResourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, application)
			}
			require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))

			deployment := resourceOf(t, stub, deploymentResourceType)
			assertApplicationVolumes(t, deployment, tt.storage)
			assertApplicationSecurity(t, deployment)

			if tt.storage == nil {
				require.Nil(t, application.Storage)
				require.Empty(t, stub.Resources(pvcResourceType))
				return
			}
			assertStorageClaim(t, resourceOf(t, stub, pvcResourceType), tt.storage)
		})
	}
}

func TestStorageIdentitySurvivesDeploymentChanges(t *testing.T) {
	configuration, err := app.ReadConfiguration("test/storage.hjson")
	require.NoError(t, err)
	configured := configuration.GetApplications()[0]
	for _, kind := range []string{"api", "worker"} {
		t.Run(kind, func(t *testing.T) {
			configured.Kind = kind
			for _, version := range []string{"1.0.0", "2.0.0"} {
				t.Run(version, func(t *testing.T) {
					configured.Version = version
					configured.Replicas++
					configured.Id = version
					stub := &test.ResourceStub{}
					run := func(ctx *pulumi.Context) error {
						return app.CreateApplication(ctx, app.ConvertApplication(configured))
					}
					require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
					claim := resourceOf(t, stub, pvcResourceType)
					metadata := test.Property(t, claim, "metadata").ObjectValue()
					require.Equal(t, "test-data", test.Property(t, metadata, "name").StringValue())
					volumes := test.Property(t, podSpec(resourceOf(t, stub, deploymentResourceType)), "volumes").ArrayValue()
					source := test.Property(t, volumes[2].ObjectValue(), "persistentVolumeClaim").ObjectValue()
					require.Equal(t, "test-data", test.Property(t, source, "claimName").StringValue())
				})
			}
		})
	}
}

func TestApplicationProtectsStorage(t *testing.T) {
	configuration, err := app.ReadConfiguration("test/storage.hjson")
	require.NoError(t, err)
	stub := &test.ResourceStub{}
	var options *pulumi.ResourceOptions
	var name string
	run := func(ctx *pulumi.Context) error {
		err := ctx.RegisterStackTransformation(func(args *pulumi.ResourceTransformationArgs) *pulumi.ResourceTransformationResult {
			if args.Type == pvcResourceType {
				name = args.Name
				options, _ = pulumi.NewResourceOptions(args.Opts...)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return app.CreateApplication(ctx, app.ConvertApplication(configuration.GetApplications()[0]))
	}
	require.NoError(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)))
	require.Equal(t, "test-data", name)
	require.NotNil(t, options)
	require.True(t, options.Protect, "data claims must require explicit unprotection before deletion or replacement")
}

func TestStorageReturnsResourceErrors(t *testing.T) {
	for _, token := range []string{pvcResourceType, deploymentResourceType} {
		t.Run(token, func(t *testing.T) {
			configuration, err := app.ReadConfiguration("test/storage.hjson")
			require.NoError(t, err)
			stub := &test.ResourceStub{}
			stub.FailResource(token)
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, app.ConvertApplication(configuration.GetApplications()[0]))
			}
			require.ErrorContains(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)), test.ErrResource.Error())
			require.NotEmpty(t, stub.Resources(token))
		})
	}
}

func TestApplicationRejectsIncompleteStorage(t *testing.T) {
	for _, tt := range []struct {
		name     string
		capacity string
		mount    string
		message  string
	}{
		{name: "missing capacity", mount: "/data", message: "storage capacity is required"},
		{name: "missing mount", capacity: "1Gi", message: "storage path"},
		{name: "relative mount", capacity: "1Gi", mount: "data", message: "storage path"},
		{name: "root mount", capacity: "1Gi", mount: "/", message: "storage path"},
		{name: "unclean mount", capacity: "1Gi", mount: "/data/../etc", message: "storage path"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			configuration, err := app.ReadConfiguration("test/storage.hjson")
			require.NoError(t, err)
			configured := configuration.GetApplications()[0]
			configured.Storage.Capacity = tt.capacity
			configured.Storage.Path = tt.mount
			stub := &test.ResourceStub{}
			run := func(ctx *pulumi.Context) error {
				return app.CreateApplication(ctx, app.ConvertApplication(configured))
			}
			require.ErrorContains(t, pulumi.RunErr(run, pulumi.WithMocks("project", "stack", stub)), tt.message)
			require.Empty(t, stub.Resources(pvcResourceType))
			require.Empty(t, stub.Resources(deploymentResourceType))
		})
	}
}

func assertStorageClaim(t *testing.T, claim resource.PropertyMap, storage *v2.Storage) {
	t.Helper()

	metadata := test.Property(t, claim, "metadata").ObjectValue()
	require.Equal(t, "test-data", test.Property(t, metadata, "name").StringValue())
	require.Equal(t, "test", test.Property(t, metadata, "namespace").StringValue())
	spec := test.Property(t, claim, "spec").ObjectValue()
	require.Equal(t, []string{"ReadWriteOnce"}, test.StringValues(test.Property(t, spec, "accessModes").ArrayValue()))
	require.Equal(t, "Filesystem", test.Property(t, spec, "volumeMode").StringValue())
	requests := test.Property(t, spec, "resources").ObjectValue()["requests"].ObjectValue()
	require.Equal(t, storage.GetCapacity(), test.Property(t, requests, "storage").StringValue())
	if storage.Profile == nil {
		require.NotContains(t, spec, resource.PropertyKey("storageClassName"))
	} else {
		require.Equal(t, storage.GetProfile(), test.Property(t, spec, "storageClassName").StringValue())
	}
}

func assertApplicationVolumes(t *testing.T, deployment resource.PropertyMap, storage *v2.Storage) {
	t.Helper()

	volumes := test.Property(t, podSpec(deployment), "volumes").ArrayValue()
	mounts := test.Property(t, container(deployment), "volumeMounts").ArrayValue()
	count := 2
	if storage != nil {
		count++
	}
	require.Len(t, volumes, count, "only opted-in apps add a data volume")
	require.Len(t, mounts, count, "only opted-in apps add a data mount")
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"name": "test", "configMap": map[string]any{"name": "test"},
	}), volumes[0].ObjectValue())
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"name": "database-secret", "secret": map[string]any{"secretName": "database-secret"},
	}), volumes[1].ObjectValue())
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"name": "test", "mountPath": "/etc/test/test.yaml", "subPath": "test.yaml",
	}), mounts[0].ObjectValue())
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"name": "database-secret", "mountPath": "/etc/secrets/database", "readOnly": true,
	}), mounts[1].ObjectValue())
	if storage != nil {
		require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
			"name": "test-data", "persistentVolumeClaim": map[string]any{"claimName": "test-data", "readOnly": false},
		}), volumes[2].ObjectValue())
		require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
			"name": "test-data", "mountPath": storage.GetPath(), "readOnly": false,
		}), mounts[2].ObjectValue())
	}
}

func assertApplicationSecurity(t *testing.T, deployment resource.PropertyMap) {
	t.Helper()

	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"runAsUser": 10001, "runAsGroup": 10001, "runAsNonRoot": true, "fsGroup": 2000,
		"fsGroupChangePolicy": "Always",
	}), test.Property(t, podSpec(deployment), "securityContext").ObjectValue())
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"readOnlyRootFilesystem": true, "allowPrivilegeEscalation": false,
		"capabilities":   map[string]any{"drop": []any{"ALL"}},
		"seccompProfile": map[string]any{"type": "RuntimeDefault"},
	}), test.Property(t, container(deployment), "securityContext").ObjectValue())
	require.Equal(t, resource.NewPropertyMapFromMap(map[string]any{
		"rollingUpdate": map[string]any{"maxSurge": 1, "maxUnavailable": 1},
	}), test.Property(t, deploymentSpec(deployment), "strategy").ObjectValue())
}

func appWithResources() *app.App {
	return &app.App{
		ID:        "1234",
		Kind:      "api",
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
