package do

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Mocks for do.
type Mocks int

func (Mocks) Call(_ pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func (Mocks) NewResource(_ pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return "", resource.PropertyMap{}, nil
}
