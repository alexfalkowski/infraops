package do

import (
	"github.com/alexfalkowski/infraops/internal/runtime"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Project for do.
type Project struct {
	Name        string
	Description string
}

// Configure for do.
func Configure(ctx *pulumi.Context) error {
	// We need a default VPC, or the first one created becomes the default one.
	r := string(digitalocean.RegionFRA1)
	n := "default-" + r
	args := &digitalocean.VpcArgs{
		Name:        pulumi.String(n),
		Region:      digitalocean.RegionFRA1,
		Description: pulumi.String("The default vpc for " + r),
	}

	_, err := digitalocean.NewVpc(ctx, n, args)

	return err
}

// CreateProject for do.
func CreateProject(ctx *pulumi.Context, project *Project) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	v, err := createVPC(ctx, project)
	runtime.Must(err)

	c, err := createCluster(ctx, v, project)
	runtime.Must(err)

	args := &digitalocean.ProjectArgs{
		Name:        pulumi.String(project.Name),
		Description: pulumi.String(project.Description),
		Environment: pulumi.String("Production"),
		Purpose:     pulumi.String("Service or API"),
		IsDefault:   pulumi.Bool(false),
		Resources: pulumi.StringArray{
			c.ClusterUrn,
		},
	}

	_, err = digitalocean.NewProject(ctx, project.Name, args)
	runtime.Must(err)

	return
}

func createVPC(ctx *pulumi.Context, p *Project) (*digitalocean.Vpc, error) {
	args := &digitalocean.VpcArgs{
		Name:        pulumi.String(p.Name),
		Region:      digitalocean.RegionFRA1,
		Description: pulumi.String(p.Description),
	}

	return digitalocean.NewVpc(ctx, p.Name, args)
}

func createCluster(ctx *pulumi.Context, v *digitalocean.Vpc, p *Project) (*digitalocean.KubernetesCluster, error) {
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
			Size:      digitalocean.DropletSlugDropletS2VCPU4GB,
		},
		Region:  pulumi.String(digitalocean.RegionFRA1),
		Version: pulumi.String("1.32.1-do.0"),
		VpcUuid: v.ID(),
	}

	return digitalocean.NewKubernetesCluster(ctx, p.Name, args)
}
