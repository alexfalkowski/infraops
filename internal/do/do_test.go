package do_test

import (
	"testing"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/do"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

const (
	kubernetesClusterResourceType = "digitalocean:index/kubernetesCluster:KubernetesCluster"
	vpcResourceType               = "digitalocean:index/vpc:Vpc"
)

func TestCreateCluster(t *testing.T) {
	for _, tt := range []struct {
		name     string
		resource string
		size     string
	}{
		{name: "small", resource: "small", size: "s-2vcpu-4gb"},
		{name: "medium", resource: "medium", size: "s-4vcpu-8gb"},
		{name: "unknown uses small", resource: "unknown", size: "s-2vcpu-4gb"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &test.ResourceStub{}
			cluster := do.ConvertCluster(&v2.Cluster{
				Name:        "test",
				Description: "test",
				Resource:    tt.resource,
			})
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				require.NoError(t, do.CreateCluster(ctx, cluster))

				return nil
			}, pulumi.WithMocks("project", "stack", stub))
			require.NoError(t, err)
			requireNodePoolSize(t, stub.Resources(kubernetesClusterResourceType), tt.size)
		})
	}

	stub := &test.ResourceStub{}
	stub.FailResource(vpcResourceType)
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		p := &do.Cluster{Name: "test", Description: "test"}
		require.NoError(t, do.CreateCluster(ctx, p))

		return nil
	}, pulumi.WithMocks("project", "stack", stub))
	require.Error(t, err)
	require.Len(t, stub.Resources(vpcResourceType), 1)
}

func TestCreateClusterReturnsKubernetesClusterError(t *testing.T) {
	stub := &test.ResourceStub{}
	stub.FailResource(kubernetesClusterResourceType)

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		p := &do.Cluster{Name: "test", Description: "test"}
		return do.CreateCluster(ctx, p)
	}, pulumi.WithMocks("project", "stack", stub))
	require.Error(t, err)
	require.Len(t, stub.Resources(vpcResourceType), 1)
	require.Len(t, stub.Resources(kubernetesClusterResourceType), 1)
}

func requireNodePoolSize(t *testing.T, clusters []resource.PropertyMap, size string) {
	t.Helper()

	require.Len(t, clusters, 1)

	nodePool := test.Property(t, clusters[0], "nodePool").ObjectValue()
	require.Equal(t, size, test.Property(t, nodePool, "size").StringValue())
}
