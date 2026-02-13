package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// prefix maps DNS record types to a stable name prefix used for Pulumi resource naming.
var prefix = map[string]string{
	"A":    "ipv4",
	"AAAA": "ipv6",
}

// BalancerZone describes a Cloudflare-managed zone where multiple hostnames are created as
// proxied A and AAAA records pointing to provided IPv4/IPv6 addresses.
type BalancerZone struct {
	// Name is the Pulumi resource name prefix used when creating Cloudflare resources for this zone.
	Name string
	// Domain is the apex domain to create/manage as a Cloudflare zone (for example "example.com").
	Domain string
	// IPV4 is the IPv4 address used for A records.
	IPV4 string
	// IPV6 is the IPv6 address used for AAAA records.
	IPV6 string
	// RecordNames are subdomain labels (for example ["api", "app"]) used to create records under Domain.
	RecordNames []string
}

// ConvertBalancerZone converts a protobuf v2.BalancerZone into the internal BalancerZone model.
func ConvertBalancerZone(z *v2.BalancerZone) *BalancerZone {
	return &BalancerZone{
		Name:        z.GetName(),
		Domain:      z.GetDomain(),
		IPV4:        z.GetIpv4(),
		IPV6:        z.GetIpv6(),
		RecordNames: z.GetRecordNames(),
	}
}

// CreateBalancerZone provisions a Cloudflare zone and creates proxied A and AAAA records.
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

// record creates a proxied DNS record of kind (for example "A" or "AAAA") with ip as content.
//
// The Pulumi resource name is prefixed (for example "ipv4.<name>") to keep the A and AAAA records
// distinct while remaining stable across updates.
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
