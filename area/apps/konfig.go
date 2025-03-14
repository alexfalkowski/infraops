package main

import (
	"github.com/alexfalkowski/infraops/internal/app"
)

func init() {
	RegisterApplication(&app.App{
		ID:        "1115c470-ccc9-4daf-8459-ef1e19c40afe",
		Name:      "konfig",
		Namespace: "lean",
		Domain:    "lean-thoughts.com",
		Version:   "1.487.0",
		Resources: &app.Resources{
			CPU:     &app.Range{Min: "250m", Max: "500m"},
			Memory:  &app.Range{Min: "128Mi", Max: "256Mi"},
			Storage: &app.Range{Min: "1Gi", Max: "2Gi"},
		},
		Secrets: app.Secrets{"konfig", "otlp", "gh"},
	})
}
