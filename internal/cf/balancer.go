package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// BalancerZone for cf.
type BalancerZone struct {
	Name        string
	Domain      string
	IPV4        string
	IPV6        string
	RecordNames []string
}

// ConvertBalancerZone converts a v2.BalancerZone to a BalancerZone.
func ConvertBalancerZone(z *v2.BalancerZone) *BalancerZone {
	return &BalancerZone{
		Name:        z.GetName(),
		Domain:      z.GetDomain(),
		IPV4:        z.GetIpv4(),
		IPV6:        z.GetIpv6(),
		RecordNames: z.GetRecordNames(),
	}
}

// CreateBalancerZone for cf.
func CreateBalancerZone(ctx *pulumi.Context, zone *BalancerZone) error {
	z, err := createZone(ctx, zone.Name, zone.Domain, "full")
	if err != nil {
		return err
	}

	for _, n := range zone.RecordNames {
		name := fmt.Sprintf("%s.%s", n, zone.Domain)

		ipv4 := &cloudflare.RecordArgs{
			Type:    pulumi.String("A"),
			Name:    pulumi.String(name),
			Content: pulumi.String(zone.IPV4),
			ZoneId:  z.ID(),
			Proxied: inputs.Yes,
			Ttl:     inputs.Automatic,
		}

		if _, err := cloudflare.NewRecord(ctx, fmt.Sprintf("ipv4.%s", name), ipv4); err != nil {
			return err
		}

		ipv6 := &cloudflare.RecordArgs{
			Type:    pulumi.String("AAAA"),
			Name:    pulumi.String(name),
			Content: pulumi.String(zone.IPV6),
			ZoneId:  z.ID(),
			Proxied: inputs.Yes,
			Ttl:     inputs.Automatic,
		}

		if _, err := cloudflare.NewRecord(ctx, fmt.Sprintf("ipv6.%s", name), ipv6); err != nil {
			return err
		}
	}

	return nil
}
