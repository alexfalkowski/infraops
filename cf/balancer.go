package cf

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// BalancerZone for cf.
type BalancerZone struct {
	Name        string
	Domain      string
	IP          string
	RecordNames []string
}

// CreateBalancerZone for cf.
func CreateBalancerZone(ctx *pulumi.Context, zone *BalancerZone) error {
	args := &cloudflare.ZoneArgs{
		AccountId: account,
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(zone.Domain),
	}

	z, err := cloudflare.NewZone(ctx, zone.Name, args)
	if err != nil {
		return err
	}

	if err := settings(ctx, zone.Name, "full", z); err != nil {
		return err
	}

	if err := dnssec(ctx, zone.Name, z); err != nil {
		return err
	}

	for _, n := range zone.RecordNames {
		name := fmt.Sprintf("%s.%s", n, zone.Domain)

		r := &cloudflare.RecordArgs{
			Type:    pulumi.String("A"),
			Name:    pulumi.String(name),
			Content: pulumi.String(zone.IP),
			ZoneId:  z.ID(),
			Proxied: pulumi.Bool(true),
			Ttl:     pulumi.Int(1),
		}

		_, err := cloudflare.NewRecord(ctx, name, r)
		if err != nil {
			return err
		}
	}

	return nil
}
