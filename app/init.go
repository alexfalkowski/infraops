package app

import (
	"strings"
)

func initConfig(app *App) (string, error) {
	cfg, err := readFile("init.yaml")
	if err != nil {
		return "", err
	}

	cfg = strings.ReplaceAll(cfg, "<app>", app.Name)
	cfg = strings.ReplaceAll(cfg, "<ver>", app.ConfigVersion)
	cfg = strings.ReplaceAll(cfg, "<ua>", app.Name+"/"+app.Version)

	return cfg, nil
}
