package app

import (
	"os"
	"path/filepath"
	"strings"

	cv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createConfigMap(ctx *pulumi.Context, app *App) error {
	var (
		args *cv1.ConfigMapArgs
		err  error
	)

	if app.ConfigVersion != "" {
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

	m := metadata(app)
	m.Name = pulumi.String("konfig")

	args := &cv1.ConfigMapArgs{
		Metadata: m,
		Data:     pulumi.StringMap{configFile("konfig"): pulumi.String(d)},
	}

	return args, nil
}

func initConfig(app *App) (string, error) {
	cfg, err := readFile("init.yaml")
	if err != nil {
		return "", err
	}

	on := []string{
		"<app>", app.Name,
		"<ver>", app.ConfigVersion,
	}
	r := strings.NewReplacer(on...)

	return r.Replace(cfg), nil
}

func configMap(app *App) (*cv1.ConfigMapArgs, error) {
	d, err := readFile(app.Name + ".yaml")
	if err != nil {
		return nil, err
	}

	args := &cv1.ConfigMapArgs{
		Metadata: metadata(app),
		Data:     pulumi.StringMap{configFile(app.Name): pulumi.String(d)},
	}

	return args, nil
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
