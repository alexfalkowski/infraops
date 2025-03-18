package main

import "github.com/alexfalkowski/infraops/internal/app"

func init() {
	RegisterApplication(&app.App{
		ID:            "28c679dc-5924-47e8-ac48-73cd842ba5cd",
		Name:          "standort",
		Namespace:     "lean",
		Domain:        "lean-thoughts.com",
		InitVersion:   "0.212.0",
		Version:       "2.362.0",
		ConfigVersion: "1.16.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "125m", Max: "250m"},
			Memory:  &app.Range{Min: "64Mi", Max: "128Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp"},
	})
}
