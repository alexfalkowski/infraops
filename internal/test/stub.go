package test

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Stub is a pulumi.Mocks implementation that returns inputs unchanged and never errors.
type Stub struct{}

// Call returns the call arguments as the result without modification.
func (*Stub) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

// NewResource returns the provided inputs as the resource state without modification.
func (*Stub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", args.Inputs, nil
}

// ErrStub is a pulumi.Mocks implementation that always returns an error.
//
// It is useful for exercising error-handling paths in Pulumi programs.
type ErrStub struct{}

// Call returns an error for every mocked provider call.
//
//nolint:err113
func (*ErrStub) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, errors.New("bad call")
}

// NewResource returns an error for every mocked resource creation.
//
//nolint:err113
func (*ErrStub) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", args.Inputs, errors.New("bad resource")
}
