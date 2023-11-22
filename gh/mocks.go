package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type mocks int

func (mocks) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func (mocks) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", resource.PropertyMap{}, nil
}
