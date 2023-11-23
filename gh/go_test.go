package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createGoHealth(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestService(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createGoService(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}
