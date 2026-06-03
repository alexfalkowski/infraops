package app

import (
	"os"
	"path/filepath"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createConfigMap(ctx *pulumi.Context, app *App) error {
	if app.IsExternal() {
		return nil
	}

	args, err := configMap(app)
	if err != nil {
		return err
	}

	_, err = cv1.NewConfigMap(ctx, app.Name, args)
	return err
}

func configMap(app *App) (*cv1.ConfigMapArgs, error) {
	configData, err := readFile(app.Namespace, app.Name+".yaml")
	if err != nil {
		return nil, err
	}

	args := &cv1.ConfigMapArgs{
		Metadata: metadata(app),
		Data:     pulumi.StringMap{configFile(app.Name): pulumi.String(configData)},
	}
	return args, nil
}

func configMatchingFilePath(name string) string {
	return configFilePath(name, name)
}

func configFilePath(path, file string) string {
	return configPath(path) + "/" + configFile(file)
}

func configPath(name string) string {
	return "/etc/" + name
}

func configFile(name string) string {
	return name + ".yaml"
}

func readFile(ns, file string) (string, error) {
	path, err := filePath(ns, file)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepath.Clean(path))
	return string(data), err
}

func filePath(ns, file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, ns, file), nil
}
