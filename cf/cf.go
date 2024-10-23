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

func settings(ctx *pulumi.Context, name string, cz *cloudflare.Zone) error {
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

	_, err := cloudflare.NewZoneSettingsOverride(ctx, name, zso)

	return err
}
