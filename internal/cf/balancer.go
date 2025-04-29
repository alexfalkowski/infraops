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
	IP          string
	RecordNames []string
}

// ConvertBalancerZone converts a v2.BalancerZone to a BalancerZone.
func ConvertBalancerZone(z *v2.BalancerZone) *BalancerZone {
	return &BalancerZone{
		Name:        z.GetName(),
		Domain:      z.GetDomain(),
		IP:          z.GetIp(),
		RecordNames: z.GetRecordNames(),
	}
}

// CreateBalancerZone for cf.
func CreateBalancerZone(ctx *pulumi.Context, zone *BalancerZone) error {
	z, err := newZone(ctx, zone.Name, zone.Domain, "full")
	if err != nil {
		return err
	}

	for _, n := range zone.RecordNames {
		name := fmt.Sprintf("%s.%s", n, zone.Domain)

		r := &cloudflare.RecordArgs{
			Type:    pulumi.String("A"),
			Name:    pulumi.String(name),
			Content: pulumi.String(zone.IP),
			ZoneId:  z.ID(),
			Proxied: inputs.Yes,
			Ttl:     inputs.Automatic,
		}

		if _, err := cloudflare.NewRecord(ctx, name, r); err != nil {
			return err
		}
	}

	return nil
}
