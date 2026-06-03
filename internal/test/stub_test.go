package test_test

import (
	"testing"

	infraopstest "github.com/alexfalkowski/infraops/v2/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

const resourceType = "test:index/resource:Resource"

func TestResourceStubRecordsResources(t *testing.T) {
	stub := &infraopstest.ResourceStub{}
	inputs := resource.PropertyMap{
		resource.PropertyKey("name"): resource.NewStringProperty("test"),
	}

	_, outputs, err := stub.NewResource(pulumi.MockResourceArgs{TypeToken: resourceType, Inputs: inputs})
	require.NoError(t, err)
	require.Equal(t, inputs, outputs)
	require.Equal(t, []resource.PropertyMap{inputs}, stub.Resources(resourceType))
}

func TestResourceStubFailsResource(t *testing.T) {
	stub := &infraopstest.ResourceStub{}
	stub.FailResource(resourceType)

	_, _, err := stub.NewResource(pulumi.MockResourceArgs{TypeToken: resourceType, Inputs: resource.PropertyMap{}})
	require.ErrorIs(t, err, infraopstest.ErrResource)
	require.Len(t, stub.Resources(resourceType), 1)
}

func TestResourceStubFailsAllResources(t *testing.T) {
	stub := &infraopstest.ResourceStub{}
	stub.FailAllResources()

	_, _, err := stub.NewResource(pulumi.MockResourceArgs{TypeToken: resourceType, Inputs: resource.PropertyMap{}})
	require.ErrorIs(t, err, infraopstest.ErrResource)
	require.Len(t, stub.Resources(resourceType), 1)
}

func TestResourceStubFailsResourceAtOccurrence(t *testing.T) {
	stub := &infraopstest.ResourceStub{}
	stub.FailResourceAt(resourceType, 2)

	_, _, err := stub.NewResource(pulumi.MockResourceArgs{TypeToken: resourceType, Inputs: resource.PropertyMap{}})
	require.NoError(t, err)

	_, _, err = stub.NewResource(pulumi.MockResourceArgs{TypeToken: resourceType, Inputs: resource.PropertyMap{}})
	require.ErrorIs(t, err, infraopstest.ErrResource)
	require.Len(t, stub.Resources(resourceType), 2)
}
