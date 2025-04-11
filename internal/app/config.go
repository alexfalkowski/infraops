package app

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createConfigMap(ctx *pulumi.Context, app *App) error {
	if !app.IsInternal() {
		return nil
	}

	var (
		args *cv1.ConfigMapArgs
		err  error
	)

	if app.HasConfigVersion() {
		args, err = initConfigMap(app)
	} else {
		args, err = configMap(app)
	}

	if err != nil {
		return err
	}

	_, err = cv1.NewConfigMap(ctx, app.Name, args)

	return err
}

func initConfigMap(app *App) (*cv1.ConfigMapArgs, error) {
	d, err := initConfig(app)
	if err != nil {
		return nil, err
	}

	name := initName(app)

	m := metadata(app)
	m.Name = pulumi.String(name)

	args := &cv1.ConfigMapArgs{
		Metadata: m,
		Data:     pulumi.StringMap{configFile(name): pulumi.String(d)},
	}

	return args, nil
}

func initName(app *App) string {
	return app.Name + "-init"
}

func initConfig(app *App) (string, error) {
	p, err := filePath(app.Namespace, "init.yaml")
	if err != nil {
		return "", err
	}

	t, err := template.ParseFiles(p)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer

	if err := t.Execute(&b, app); err != nil {
		return "", err
	}

	return b.String(), nil
}

func configMap(app *App) (*cv1.ConfigMapArgs, error) {
	d, err := readFile(app.Namespace, app.Name+".yaml")
	if err != nil {
		return nil, err
	}

	args := &cv1.ConfigMapArgs{
		Metadata: metadata(app),
		Data:     pulumi.StringMap{configFile(app.Name): pulumi.String(d)},
	}

	return args, nil
}

func configMatchingFilePath(name string) pulumi.String {
	return configFilePath(name, name)
}

func configFilePath(path, file string) pulumi.String {
	return pulumi.String(configPath(path) + "/" + configFile(file))
}

func configPath(name string) string {
	return "/etc/" + name
}

func configFile(name string) string {
	return name + ".yaml"
}

func readFile(ns, file string) (string, error) {
	p, err := filePath(ns, file)
	if err != nil {
		return "", err
	}

	d, err := os.ReadFile(filepath.Clean(p))

	return string(d), err
}

func filePath(ns, file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, ns, file), nil
}
