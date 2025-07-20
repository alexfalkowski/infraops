package main

import (
	"flag"
	"log"
	"os"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
)

var configs = map[string]config.Config{
	"apps": &v2.Kubernetes{},
	"cf":   &v2.Cloudflare{},
	"do":   &v2.DigitalOcean{},
	"gh":   &v2.Github{},
}

func main() {
	var (
		kind string
		path string
	)

	set := flag.NewFlagSet("format", flag.ContinueOnError)
	set.StringVar(&kind, "k", "", "kind of config")
	set.StringVar(&path, "p", "", "path of the config")
	if err := set.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	cfg, ok := configs[kind]
	if !ok {
		log.Fatal("invalid kind")
	}

	if err := config.Read(path, cfg); err != nil {
		log.Fatal(err)
	}

	if err := config.Write(path, cfg); err != nil {
		log.Fatal(err)
	}
}
