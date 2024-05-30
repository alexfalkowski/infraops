package cf

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Zone for cf.
type Zone struct {
	Name        string
	Domain      string
	Balancer    string
	RecordNames []string
}

// CreateZone for cf.
func CreateZone(ctx *pulumi.Context, zone *Zone) error {
	a, err := account(ctx)
	if err != nil {
		return err
	}

	args := &cloudflare.ZoneArgs{
		AccountId: a.ID(),
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(zone.Domain),
	}

	z, err := cloudflare.NewZone(ctx, zone.Name, args)
	if err != nil {
		return err
	}

	for _, n := range zone.RecordNames {
		args := &cloudflare.RecordArgs{
			Type:    pulumi.String("A"),
			Name:    pulumi.String(n),
			Value:   pulumi.String(zone.Balancer),
			ZoneId:  z.ID(),
			Proxied: pulumi.Bool(true),
			Ttl:     pulumi.Int(1),
		}

		_, err := cloudflare.NewRecord(ctx, n, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func account(ctx *pulumi.Context) (*cloudflare.Account, error) {
	args := &cloudflare.AccountArgs{
		Name: pulumi.String("main account"),
	}

	return cloudflare.NewAccount(ctx, "main", args)
}
