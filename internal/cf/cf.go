package cf

import (
	"fmt"
	"os"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/alexfalkowski/infraops/v2/internal/inputs"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// account is the Cloudflare account identifier sourced from the CLOUDFLARE_ACCOUNT_ID
// environment variable.
//
// It is intentionally a package-level Pulumi input to avoid threading the value through
// every resource constructor.
var account = pulumi.String(os.Getenv("CLOUDFLARE_ACCOUNT_ID"))

// ReadConfiguration reads the Cloudflare area configuration from path.
//
// The file is expected to be HJSON matching the v2.Cloudflare protobuf schema.
func ReadConfiguration(path string) (*v2.Cloudflare, error) {
	var configuration v2.Cloudflare
	err := config.Read(path, &configuration)
	return &configuration, err
}

// ZoneSetting represents a single Cloudflare zone setting name/value pair.
//
// It is used to apply a consistent baseline of zone settings during provisioning.
type ZoneSetting struct {
	Name  string
	Value pulumi.String
}

// settings applies a baseline set of Cloudflare zone settings to the given zone.
//
// The resource names are derived from the provided name plus the setting identifier to keep
// them stable across updates. The ssl argument is applied to the zone "ssl" setting.
func settings(ctx *pulumi.Context, name, ssl string, cz *cloudflare.Zone) error {
	settings := []*ZoneSetting{
		{Name: "always_use_https", Value: inputs.On},
		{Name: "min_tls_version", Value: pulumi.String("1.2")},
		{Name: "cache_level", Value: pulumi.String("aggressive")},
		{Name: "http3", Value: inputs.On},
		{Name: "email_obfuscation", Value: inputs.Off},
		{Name: "h2_prioritization", Value: inputs.On},
		{Name: "ssl", Value: pulumi.String(ssl)},
	}
	for _, setting := range settings {
		args := &cloudflare.ZoneSettingArgs{
			SettingId: pulumi.String(setting.Name),
			Value:     setting.Value,
			ZoneId:    cz.ID(),
		}
		name := fmt.Sprintf("%s_%s", name, setting.Name)

		if _, err := cloudflare.NewZoneSetting(ctx, name, args); err != nil {
			return err
		}
	}

	return nil
}
