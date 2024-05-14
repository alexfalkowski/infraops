package do

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Project for do.
type Project struct {
	Name        string
	Description string
}

// CreateProject for do.
func CreateProject(ctx *pulumi.Context, project *Project) error {
	v, err := createVPC(ctx, project)
	if err != nil {
		return err
	}

	c, err := createCluster(ctx, v, project)
	if err != nil {
		return err
	}

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

	return err
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
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			Name:      pulumi.String(p.Name),
			Size:      digitalocean.DropletSlugDropletS1VCPU2GB,
			AutoScale: pulumi.Bool(false),
			MaxNodes:  pulumi.Int(1),
			MinNodes:  pulumi.Int(1),
			NodeCount: pulumi.Int(1),
		},
		Region:                        digitalocean.RegionFRA1,
		Version:                       pulumi.String("1.29.1-do.0"),
		AutoUpgrade:                   pulumi.Bool(false),
		DestroyAllAssociatedResources: pulumi.Bool(true),
		Ha:                            pulumi.Bool(false),
		Name:                          pulumi.String(p.Name),
		RegistryIntegration:           pulumi.Bool(false),
		SurgeUpgrade:                  pulumi.Bool(false),
		VpcUuid:                       v.VpcUrn,
	}

	return digitalocean.NewKubernetesCluster(ctx, p.Name, args)
}
