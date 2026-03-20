package main

import (
	"path/filepath"
	"testing"

	"github.com/alexfalkowski/infraops/v2/internal/cf"
	"github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestBalancerZones(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := cf.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			zones := config.GetBalancerZones()
			run := func(ctx *pulumi.Context) error {
				for _, zone := range zones {
					if err := cf.CreateBalancerZone(ctx, cf.ConvertBalancerZone(zone)); err != nil {
						return err
					}
				}

				return nil
			}

			require.NoError(t, runWithMocks(run, &test.Stub{}))
			require.Error(t, runWithMocks(run, &test.ErrStub{}))
		})
	}
}

func TestPageZones(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := cf.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			zones := config.GetPageZones()
			run := func(ctx *pulumi.Context) error {
				for _, zone := range zones {
					if err := cf.CreatePageZone(ctx, cf.ConvertPageZone(zone)); err != nil {
						return err
					}
				}

				return nil
			}

			require.NoError(t, runWithMocks(run, &test.Stub{}))
			require.Error(t, runWithMocks(run, &test.ErrStub{}))
		})
	}
}

func TestBuckets(t *testing.T) {
	for _, fixture := range fixtures() {
		t.Run(fixture.name, func(t *testing.T) {
			config, err := cf.ReadConfiguration(fixture.path)
			require.NoError(t, err)

			buckets := config.GetBuckets()
			run := func(ctx *pulumi.Context) error {
				for _, bucket := range buckets {
					if err := cf.CreateBucket(ctx, cf.ConvertBucket(bucket)); err != nil {
						return err
					}
				}

				return nil
			}

			require.NoError(t, runWithMocks(run, &test.Stub{}))
			if len(buckets) == 0 {
				return
			}

			require.Error(t, runWithMocks(run, &test.ErrStub{}))
		})
	}
}

type fixture struct {
	name string
	path string
}

func fixtures() []fixture {
	return []fixture{
		{name: "area", path: "cf.hjson"},
		{name: "shared", path: filepath.Join("..", "..", "internal", "test", "cf.hjson")},
	}
}

func runWithMocks(run pulumi.RunFunc, mocks pulumi.MockResourceMonitor) error {
	return pulumi.RunErr(run, pulumi.WithMocks("project", "stack", mocks))
}
