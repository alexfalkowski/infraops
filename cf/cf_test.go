package cf_test

import (
	"testing"

	"github.com/alexfalkowski/infraops/cf"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateZone(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		z := &cf.Zone{
			Name:        "test",
			Domain:      "test.com",
			RecordNames: []string{"test"},
			Balancer:    "127.0.0.1",
		}

		err := cf.CreateZone(ctx, z)
		require.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", cf.Mocks(0)))

	require.NoError(t, err)
}
