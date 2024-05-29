package app

import (
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func metadata(app *App) mv1.ObjectMetaArgs {
	return mv1.ObjectMetaArgs{
		Name:      pulumi.String(app.Name),
		Namespace: pulumi.String(app.Name),
	}
}

func labels(app *App) pulumi.StringMap {
	return pulumi.StringMap{
		"app": pulumi.String(app.Name),
	}
}

func annotations(app *App) pulumi.StringMap {
	f := pulumi.String("false")

	return pulumi.StringMap{
		"circleci.com/project-id":                pulumi.String(app.ID),
		"circleci.com/restore-version-enabled":   f,
		"circleci.com/scale-component-enabled":   f,
		"circleci.com/restart-component-enabled": f,
		"circleci.com/retry-release-enabled":     f,
		"circleci.com/promote-release-enabled":   f,
		"circleci.com/cancel-release-enabled":    f,
	}
}
