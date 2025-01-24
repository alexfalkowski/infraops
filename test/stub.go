package test

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Stub for test.
type Stub struct{}

func (*Stub) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func (*Stub) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", resource.PropertyMap{}, nil
}

// ErrStub for test.
type ErrStub struct{}

//nolint:err113
func (*ErrStub) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return nil, errors.New("bad call")
}

//nolint:err113
func (*ErrStub) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", nil, errors.New("bad resource")
}
