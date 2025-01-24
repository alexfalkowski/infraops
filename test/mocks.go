package test

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Mocks for test.
type Mocks int

func (Mocks) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func (Mocks) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", resource.PropertyMap{}, nil
}

// BadMocks for test.
type BadMocks int

//nolint:err113
func (BadMocks) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return nil, errors.New("bad call")
}

//nolint:err113
func (BadMocks) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", nil, errors.New("bad resource")
}
