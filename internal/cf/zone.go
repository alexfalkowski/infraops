package cf

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func newZone(ctx *pulumi.Context, name, domain, ssl string) (*cloudflare.Zone, error) {
	args := &cloudflare.ZoneArgs{
		Account: cloudflare.ZoneAccountArgs{Id: account},
		Name:    pulumi.String(domain),
		Type:    pulumi.String("full"),
	}

	z, err := cloudflare.NewZone(ctx, name, args)
	if err != nil {
		return nil, err
	}

	if err := settings(ctx, name, ssl, z); err != nil {
		return nil, err
	}

	return z, nil
}
