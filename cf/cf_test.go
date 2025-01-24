package cf_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/cf"
	"github.com/alexfalkowski/infraops/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateBalancerZone(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		z := &cf.BalancerZone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			IP:          "127.0.0.1",
		}

		err := cf.CreateBalancerZone(ctx, z)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.Mocks(0)))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		z := &cf.BalancerZone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			IP:          "127.0.0.1",
		}

		err := cf.CreateBalancerZone(ctx, z)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.BadMocks(0)))
	require.Error(t, err)
}

func TestCreatePagerZone(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		z := &cf.PageZone{
			Name:   "test",
			Domain: "test.com",
			Host:   "test.github.io",
		}

		err := cf.CreatePageZone(ctx, z)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.Mocks(0)))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		z := &cf.PageZone{
			Name:   "test",
			Domain: "test.com",
			Host:   "test.github.io",
		}

		err := cf.CreatePageZone(ctx, z)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", test.BadMocks(0)))
	require.Error(t, err)
}
