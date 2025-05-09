package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var prefix = map[string]string{
	"A":    "ipv4",
	"AAAA": "ipv6",
}

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

		if err := record(ctx, name, "A", zone.IPV4, z.ID()); err != nil {
			return err
		}

		if err := record(ctx, name, "AAAA", zone.IPV6, z.ID()); err != nil {
			return err
		}
	}

	return nil
}

func record(ctx *pulumi.Context, name, kind, ip string, id pulumi.IDOutput) error {
	record := &cloudflare.RecordArgs{
		Type:    pulumi.String(kind),
		Name:    pulumi.String(name),
		Content: pulumi.String(ip),
		ZoneId:  id,
		Proxied: inputs.Yes,
		Ttl:     inputs.Automatic,
	}
	_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s.%s", prefix[kind], name), record)

	return err
}
