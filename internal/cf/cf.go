package cf

import (
	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/config"
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

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Cloudflare, error) {
	var configuration v2.Cloudflare
	err := config.Read(path, &configuration)

	return &configuration, err
}

func settings(ctx *pulumi.Context, name, ssl string, cz *cloudflare.Zone) error {
	ss := cloudflare.ZoneSettingsOverrideSettingsSecurityHeaderArgs{
		Enabled:           yes,
		IncludeSubdomains: yes,
		Nosniff:           yes,
		Preload:           yes,
		MaxAge:            year,
	}

	st := &cloudflare.ZoneSettingsOverrideSettingsArgs{
		AlwaysUseHttps:   on,
		MinTlsVersion:    pulumi.String("1.2"),
		CacheLevel:       pulumi.String("aggressive"),
		Http3:            on,
		EmailObfuscation: off,
		H2Prioritization: on,
		SecurityHeader:   ss,
		Ssl:              pulumi.String(ssl),
	}

	zso := &cloudflare.ZoneSettingsOverrideArgs{
		ZoneId:   cz.ID(),
		Settings: st,
	}

	_, err := cloudflare.NewZoneSettingsOverride(ctx, name, zso)

	return err
}

func dnssec(ctx *pulumi.Context, name string, cz *cloudflare.Zone) error {
	_, err := cloudflare.NewZoneDnssec(ctx, name, &cloudflare.ZoneDnssecArgs{
		ZoneId: cz.ID(),
	})

	return err
}
