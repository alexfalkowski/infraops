package pulumi

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Stub for test.
type Stub struct{}

func (*Stub) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func (*Stub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", args.Inputs, nil
}

// ErrStub for test.
type ErrStub struct{}

//nolint:err113
func (*ErrStub) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, errors.New("bad call")
}

//nolint:err113
func (*ErrStub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", args.Inputs, errors.New("bad resource")
}
