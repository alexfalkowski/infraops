package main

import (
	"testing"

	"github.com/alexfalkowski/infraops/internal/do"
	"github.com/alexfalkowski/infraops/internal/test"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	config, err := do.ReadConfiguration("do.pbtxt")
	require.NoError(t, err)

	projects := config.GetProjects()

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, project := range projects {
			err := do.CreateProject(ctx, do.ConvertProject(project))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.Stub{}))
	require.NoError(t, err)

	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
		for _, project := range projects {
			err := do.CreateProject(ctx, do.ConvertProject(project))
			require.NoError(t, err)
		}

		return nil
	}, pulumi.WithMocks("project", "stack", &test.ErrStub{}))
	require.Error(t, err)
}
