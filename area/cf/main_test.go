package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/cf"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

//nolint:dupl
func TestBalancerZones(t *testing.T) {
	config, err := cf.ReadConfiguration("cf.hjson")
	require.NoError(t, err)

	zones := config.GetBalancerZones()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, zone := range zones {
			err := cf.CreateBalancerZone(ctx, cf.ConvertBalancerZone(zone))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, zone := range zones {
			err := cf.CreateBalancerZone(ctx, cf.ConvertBalancerZone(zone))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}

//nolint:dupl
func TestPageZones(t *testing.T) {
	config, err := cf.ReadConfiguration("cf.hjson")
	require.NoError(t, err)

	zones := config.GetPageZones()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, zone := range zones {
			err := cf.CreatePageZone(ctx, cf.ConvertPageZone(zone))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, zone := range zones {
			err := cf.CreatePageZone(ctx, cf.ConvertPageZone(zone))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}

func TestBuckets(t *testing.T) {
	config, err := cf.ReadConfiguration("cf.hjson")
	require.NoError(t, err)

	buckets := config.GetBuckets()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, bucket := range buckets {
			err := cf.CreateBucket(ctx, cf.ConvertBucket(bucket))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)
}
