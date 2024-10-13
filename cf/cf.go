package cf

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	on      = pulumi.String("on")
	off     = pulumi.String("off")
	yes     = pulumi.Bool(true)
	year    = pulumi.Int(31536000)
	account = pulumi.String("561357e2a2b66ddfeabd46e2965d2c67")
)

// Zone for cf.
type Zone struct {
	Name        string
	Domain      string
	Balancer    string
	RecordNames []string
}

// CreateZone for cf.
func CreateZone(ctx *pulumi.Context, zone *Zone) error {
	args := &cloudflare.ZoneArgs{
		AccountId: account,
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(zone.Domain),
	}

	z, err := cloudflare.NewZone(ctx, zone.Name, args)
	if err != nil {
		return err
	}

	if err := settings(ctx, zone, z); err != nil {
		return err
	}

	for _, n := range zone.RecordNames {
		r := &cloudflare.RecordArgs{
			Type:    pulumi.String("A"),
			Name:    pulumi.String(n),
			Content: pulumi.String(zone.Balancer),
			ZoneId:  z.ID(),
			Proxied: pulumi.Bool(true),
			Ttl:     pulumi.Int(1),
		}

		_, err := cloudflare.NewRecord(ctx, n, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func settings(ctx *pulumi.Context, zone *Zone, cz *cloudflare.Zone) error {
	ss := cloudflare.ZoneSettingsOverrideSettingsSecurityHeaderArgs{
		Enabled:           yes,
		IncludeSubdomains: yes,
		Nosniff:           yes,
		Preload:           yes,
		MaxAge:            year,
	}

	st := &cloudflare.ZoneSettingsOverrideSettingsArgs{
		MinTlsVersion:    pulumi.String("1.2"),
		CacheLevel:       pulumi.String("aggressive"),
		Http3:            on,
		EmailObfuscation: off,
		H2Prioritization: on,
		SecurityHeader:   ss,
	}

	zso := &cloudflare.ZoneSettingsOverrideArgs{
		ZoneId:   cz.ID(),
		Settings: st,
	}

	_, err := cloudflare.NewZoneSettingsOverride(ctx, zone.Name, zso)

	return err
}
