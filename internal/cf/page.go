package cf

import (
	"fmt"

	v1 "github.com/alexfalkowski/infraops/api/infraops/v1"
	"github.com/alexfalkowski/infraops/internal/runtime"
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PageZone for cf.
type PageZone struct {
	Name   string
	Domain string
	Host   string
}

// ConvertPageZone converts a v1.PageZone to a PageZone.
func ConvertPageZone(z *v1.PageZone) *PageZone {
	return &PageZone{
		Name:   z.GetName(),
		Domain: z.GetDomain(),
		Host:   z.GetHost(),
	}
}

// CreatePageZone for cf.
func CreatePageZone(ctx *pulumi.Context, zone *PageZone) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	args := &cloudflare.ZoneArgs{
		AccountId: account,
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(zone.Domain),
	}

	z, err := cloudflare.NewZone(ctx, zone.Name, args)
	runtime.Must(err)

	err = settings(ctx, zone.Name, "strict", z)
	runtime.Must(err)

	err = dnssec(ctx, zone.Name, z)
	runtime.Must(err)

	name := fmt.Sprintf("%s.%s", "www", zone.Domain)
	r := &cloudflare.RecordArgs{
		Type:    pulumi.String("CNAME"),
		Name:    pulumi.String(name),
		Content: pulumi.String(zone.Host),
		ZoneId:  z.ID(),
		Proxied: pulumi.Bool(true),
		Ttl:     pulumi.Int(1),
	}

	_, err = cloudflare.NewRecord(ctx, name, r)
	runtime.Must(err)

	return
}
