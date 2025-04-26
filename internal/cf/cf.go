package cf

import (
	"os"

	v2 "github.com/alexfalkowski/infraops/api/infraops/v2"
	"github.com/alexfalkowski/infraops/internal/config"
	"github.com/alexfalkowski/infraops/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var account = pulumi.String(os.Getenv("CLOUDFLARE_ACCOUNT_ID"))

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Cloudflare, error) {
	var configuration v2.Cloudflare
	err := config.Read(path, &configuration)

	return &configuration, err
}

func settings(ctx *pulumi.Context, name, ssl string, cz *cloudflare.Zone) error {
	ss := cloudflare.ZoneSettingsOverrideSettingsSecurityHeaderArgs{
		Enabled:           inputs.Yes,
		IncludeSubdomains: inputs.Yes,
		Nosniff:           inputs.Yes,
		Preload:           inputs.Yes,
		MaxAge:            pulumi.Int(31536000),
	}

	st := &cloudflare.ZoneSettingsOverrideSettingsArgs{
		AlwaysUseHttps:   inputs.On,
		MinTlsVersion:    pulumi.String("1.2"),
		CacheLevel:       pulumi.String("aggressive"),
		Http3:            inputs.On,
		EmailObfuscation: inputs.Off,
		H2Prioritization: inputs.On,
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
