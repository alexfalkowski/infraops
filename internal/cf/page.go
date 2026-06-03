package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PageZone describes a Cloudflare zone used for a static-site or pages-style CNAME target.
type PageZone struct {
	// Name is the Pulumi resource name prefix used when creating Cloudflare resources for this zone.
	Name string
	// Domain is the apex domain to create/manage as a Cloudflare zone (for example "example.com").
	Domain string
	// Host is the DNS hostname to CNAME to (for example "<owner>.github.io" or "<project>.pages.dev").
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

// CreatePageZone provisions a Cloudflare zone for a static-site CNAME target.
//
// It applies the shared zone-settings baseline with "strict" SSL mode, then creates a proxied
// CNAME from "www.<domain>" to zone.Host using Cloudflare's automatic TTL.
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
