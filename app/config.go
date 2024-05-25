package app

import (
	"os"
	"path/filepath"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createConfigMap(ctx *pulumi.Context, app *App) error {
	if app.ConfigVersion != "" {
		d, err := generateInitConfig(app)
		if err != nil {
			return err
		}

		args := &cv1.ConfigMapArgs{
			Metadata: metadata(app),
			Data:     pulumi.StringMap{configFile("konfig"): pulumi.String(d)},
		}
		_, err = cv1.NewConfigMap(ctx, app.Name, args)

		return err
	}

	d, err := readFile(app.Name + ".yaml")
	if err != nil {
		return err
	}

	args := &cv1.ConfigMapArgs{
		Metadata: metadata(app),
		Data:     pulumi.StringMap{configFile(app.Name): pulumi.String(d)},
	}
	_, err = cv1.NewConfigMap(ctx, app.Name, args)

	return err
}

func configFullFilePath(name string) pulumi.String {
	return pulumi.String(configPath(name) + "/" + configFile(name))
}

func configPath(name string) string {
	return "/etc/" + name
}

func configFile(name string) string {
	return name + ".yaml"
}

func readFile(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	p := filepath.Clean(filepath.Join(wd, file))
	d, err := os.ReadFile(p)

	return string(d), err
}
