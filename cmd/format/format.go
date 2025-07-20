package main

import (
	"flag"
	"log"
	"os"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
)

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

	var cfg config.Config

	switch kind {
	case "apps":
		cfg = &v2.Kubernetes{}
	case "cf":
		cfg = &v2.Cloudflare{}
	case "do":
		cfg = &v2.DigitalOcean{}
	case "gh":
		cfg = &v2.Github{}
	}

	if cfg == nil {
		log.Fatal("invalid kind")
	}

	if err := config.Read(path, cfg); err != nil {
		log.Fatal(err)
	}

	if err := config.Write(path, cfg); err != nil {
		log.Fatal(err)
	}
}
