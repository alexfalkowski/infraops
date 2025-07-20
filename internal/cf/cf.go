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

var account = pulumi.String(os.Getenv("CLOUDFLARE_ACCOUNT_ID"))

// ReadConfiguration reads a file and populates a configuration.
func ReadConfiguration(path string) (*v2.Cloudflare, error) {
	var configuration v2.Cloudflare
	err := config.Read(path, &configuration)
	return &configuration, err
}

// ZoneSetting is a name value.
type ZoneSetting struct {
	Name  string
	Value pulumi.String
}

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
