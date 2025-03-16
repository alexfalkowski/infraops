package main

import "github.com/alexfalkowski/infraops/internal/app"

func init() {
	RegisterApplication(&app.App{
		ID:            "98968ca0-4ada-4856-8547-210f92b838ea",
		Name:          "bezeichner",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "0.211.0",
		Version:       "1.367.0",
		ConfigVersion: "1.15.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp"},
	})
}
