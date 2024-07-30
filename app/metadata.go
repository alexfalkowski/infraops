package app

import (
	mv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func metadata(app *App, ms ...pulumi.StringMap) mv1.ObjectMetaArgs {
	ms = append(ms, recommendedLabels(app))

	return mv1.ObjectMetaArgs{
		Name:      pulumi.String(app.Name),
		Namespace: pulumi.String(app.Namespace),
		Labels:    merge(ms...),
	}
}

func recommendedLabels(app *App) pulumi.StringMap {
	return pulumi.StringMap{
		"app.kubernetes.io/name":    pulumi.String(app.Name),
		"app.kubernetes.io/version": pulumi.String(app.Version),
	}
}

func matchLabels(app *App) pulumi.StringMap {
	return pulumi.StringMap{
		"app": pulumi.String(app.Name),
	}
}

func deploymentLabels(app *App) pulumi.StringMap {
	return pulumi.StringMap{
		"circleci.com/component-name": pulumi.String(app.Name),
		"circleci.com/version":        pulumi.String(app.Version),
	}
}

func deploymentAnnotations(app *App) pulumi.StringMap {
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

func merge(ms ...pulumi.StringMap) pulumi.StringMap {
	fm := pulumi.StringMap{}

	for _, m := range ms {
		for k, v := range m {
			fm[k] = v
		}
	}

	return fm
}
