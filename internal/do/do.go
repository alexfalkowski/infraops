package do

import (
	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var sizes = map[string]digitalocean.DropletSlug{
	"small":  digitalocean.DropletSlugDropletS2VCPU4GB,
	"medium": digitalocean.DropletSlugDropletS4VCPU8GB,
}

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.DigitalOcean, error) {
	var configuration v2.DigitalOcean
	err := config.Read(path, &configuration)
	return &configuration, err
}

// Cluster for do.
type Cluster struct {
	Name        string
	Description string
	Resource    string
}

// Size of the cluster.
func (c *Cluster) Size() digitalocean.DropletSlug {
	if s, ok := sizes[c.Resource]; ok {
		return s
	}
	return digitalocean.DropletSlugDropletS2VCPU4GB
}

// ConvertCluster converts a v2.Cluster to a Cluster.
func ConvertCluster(cluster *v2.Cluster) *Cluster {
	return &Cluster{
		Name:        cluster.GetName(),
		Description: cluster.GetDescription(),
		Resource:    cluster.GetResource(),
	}
}

// CreateCluster for do.
func CreateCluster(ctx *pulumi.Context, cluster *Cluster) (err error) {
	v, err := createVPC(ctx, cluster)
	if err != nil {
		return err
	}

	_, err = createCluster(ctx, v, cluster)
	return err
}

func createVPC(ctx *pulumi.Context, p *Cluster) (*digitalocean.Vpc, error) {
	args := &digitalocean.VpcArgs{
		Name:        pulumi.String(p.Name),
		Region:      digitalocean.RegionFRA1,
		Description: pulumi.String(p.Description),
	}
	return digitalocean.NewVpc(ctx, p.Name, args)
}

func createCluster(ctx *pulumi.Context, vpc *digitalocean.Vpc, cluster *Cluster) (*digitalocean.KubernetesCluster, error) {
	args := &digitalocean.KubernetesClusterArgs{
		MaintenancePolicy: &digitalocean.KubernetesClusterMaintenancePolicyArgs{
			Day:       pulumi.String("any"),
			StartTime: pulumi.String("23:00"),
		},
		Name:                          pulumi.String(cluster.Name),
		DestroyAllAssociatedResources: inputs.Yes,
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			NodeCount: pulumi.Int(2),
			Name:      pulumi.String(cluster.Name),
			Labels:    pulumi.StringMap{"name": pulumi.String(cluster.Name)},
			Size:      cluster.Size(),
		},
		Region:  pulumi.String(digitalocean.RegionFRA1),
		Version: pulumi.String("1.34.1-do.3"),
		VpcUuid: vpc.ID(),
	}
	return digitalocean.NewKubernetesCluster(ctx, cluster.Name, args)
}
