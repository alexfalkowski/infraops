package cf

import (
	"fmt"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
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
	args := &cloudflare.ZoneArgs{
		AccountId: account,
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(zone.Domain),
	}

	z, err := cloudflare.NewZone(ctx, zone.Name, args)
	if err != nil {
		return err
	}

	if err := settings(ctx, zone.Name, "strict", z); err != nil {
		return err
	}

	if err := dnssec(ctx, zone.Name, z); err != nil {
		return err
	}

	name := fmt.Sprintf("%s.%s", "www", zone.Domain)
	r := &cloudflare.RecordArgs{
		Type:    pulumi.String("CNAME"),
		Name:    pulumi.String(name),
		Content: pulumi.String(zone.Host),
		ZoneId:  z.ID(),
		Proxied: inputs.Yes,
		Ttl:     inputs.One,
	}

	_, err = cloudflare.NewRecord(ctx, name, r)

	return err
}
