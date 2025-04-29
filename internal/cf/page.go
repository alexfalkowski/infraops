package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PageZone for cf.
type PageZone struct {
	Name   string
	Domain string
	Host   string
}

// ConvertPageZone converts a v2.PageZone to a PageZone.
func ConvertPageZone(z *v2.PageZone) *PageZone {
	return &PageZone{
		Name:   z.GetName(),
		Domain: z.GetDomain(),
		Host:   z.GetHost(),
	}
}

// CreatePageZone for cf.
func CreatePageZone(ctx *pulumi.Context, zone *PageZone) error {
	z, err := newZone(ctx, zone.Name, zone.Domain, "strict")
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
