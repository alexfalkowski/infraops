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

	args := &digitalocean.ProjectArgs{
		Name:        pulumi.String(project.Name),
		Description: pulumi.String(project.Description),
		Environment: pulumi.String("Production"),
		Purpose:     pulumi.String("Service or API"),
		IsDefault:   pulumi.Bool(false),
		Resources: pulumi.StringArray{
			v.VpcUrn,
		},
	}
	_, err = digitalocean.NewProject(ctx, project.Name, args)

	return err
}

func createVPC(ctx *pulumi.Context, p *Project) (*digitalocean.Vpc, error) {
	args := &digitalocean.VpcArgs{
		Region:      pulumi.String("fra1"),
		Description: pulumi.String(p.Description),
	}

	return digitalocean.NewVpc(ctx, p.Name, args)
}
