package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createTemplate(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestStatus(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createStatus(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestStandort(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createStandort(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestAuth(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createAuth(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestKonfig(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createKonfig(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}

func TestMigrieren(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		err := createMigrieren(ctx)
		assert.NoError(t, err)

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	assert.NoError(t, err)
}
