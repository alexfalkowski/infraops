package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PageZone describes a Cloudflare Pages-backed zone and the DNS configuration required for it.
type PageZone struct {
	// Name is the Pulumi resource name prefix used when creating Cloudflare resources for this zone.
	Name string
	// Domain is the apex domain to create/manage as a Cloudflare zone (for example "example.com").
	Domain string
	// Host is the Cloudflare Pages hostname to CNAME to (for example "<project>.pages.dev").
	Host string
}

// ConvertPageZone converts a protobuf v2.PageZone into the internal PageZone model.
func ConvertPageZone(z *v2.PageZone) *PageZone {
	return &PageZone{
		Name:   z.GetName(),
		Domain: z.GetDomain(),
		Host:   z.GetHost(),
	}
}

// CreatePageZone provisions a Cloudflare zone for a Pages site and creates a proxied CNAME record.
//
// It creates the zone with "strict" SSL mode and then creates a CNAME from "www.<domain>"
// to zone.Host. The DNS record is proxied and uses Cloudflare's automatic TTL.
func CreatePageZone(ctx *pulumi.Context, zone *PageZone) error {
	z, err := createZone(ctx, zone.Name, zone.Domain, "strict")
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s.%s", "www", zone.Domain)
	r := &cloudflare.RecordArgs{
		Type:    pulumi.String("CNAME"),
		Name:    pulumi.String(name),
		Content: pulumi.String(zone.Host),
		ZoneId:  z.ID(),
		Proxied: inputs.Yes,
		Ttl:     inputs.Automatic,
	}
	_, err = cloudflare.NewRecord(ctx, name, r)

	return err
}
