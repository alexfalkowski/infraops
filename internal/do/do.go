package do

import (
	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/config"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

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
}

// ConvertCluster converts a v2.Cluster to a Cluster.
func ConvertCluster(cluster *v2.Cluster) *Cluster {
	return &Cluster{
		Name:        cluster.GetName(),
		Description: cluster.GetDescription(),
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

func createCluster(ctx *pulumi.Context, v *digitalocean.Vpc, p *Cluster) (*digitalocean.KubernetesCluster, error) {
	args := &digitalocean.KubernetesClusterArgs{
		MaintenancePolicy: &digitalocean.KubernetesClusterMaintenancePolicyArgs{
			Day:       pulumi.String("any"),
			StartTime: pulumi.String("23:00"),
		},
		Name:                          pulumi.String(p.Name),
		DestroyAllAssociatedResources: pulumi.Bool(true),
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			NodeCount: pulumi.Int(2),
			Name:      pulumi.String(p.Name),
			Size:      digitalocean.DropletSlugDropletS4VCPU8GB,
		},
		Region:  pulumi.String(digitalocean.RegionFRA1),
		Version: pulumi.String("1.32.2-do.0"),
		VpcUuid: v.ID(),
	}

	return digitalocean.NewKubernetesCluster(ctx, p.Name, args)
}
