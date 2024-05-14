package cf

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateZone for cf.
func CreateZone(ctx *pulumi.Context, name string) error {
	a, err := account(ctx)
	if err != nil {
		return err
	}

	args := &cloudflare.ZoneArgs{
		AccountId: a.ID(),
		Plan:      pulumi.String("free"),
		Zone:      pulumi.String(name + ".com"),
	}
	_, err = cloudflare.NewZone(ctx, name, args)

	return err
}

func account(ctx *pulumi.Context) (*cloudflare.Account, error) {
	args := &cloudflare.AccountArgs{
		Name: pulumi.String("main account"),
	}

	return cloudflare.NewAccount(ctx, "main", args)
}
