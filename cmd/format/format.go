// Command format normalizes/rewrites an area configuration file in a consistent HJSON form.
//
// It reads a configuration into the corresponding protobuf message and writes it back out
// using the repository's canonical HJSON formatting rules.
//
// Usage:
//
//	format -k <apps|cf|do|gh> [-p <path>]
//
// By default, the path is `area/<kind>/<kind>.hjson`.
package main

import (
	"flag"
	"fmt"
	"os"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/alexfalkowski/infraops/v2/internal/config"
	"github.com/alexfalkowski/infraops/v2/internal/log"
)

// configs maps a supported config kind (as provided by the `-k` flag) to the protobuf
// message type used to decode/encode that configuration.
var configs = map[string]config.Config{
	"apps": &v2.Kubernetes{},
	"cf":   &v2.Cloudflare{},
	"do":   &v2.DigitalOcean{},
	"gh":   &v2.Github{},
}

// run parses CLI flags, loads the selected configuration, and writes it back out.
//
// The `-k` flag selects which protobuf message schema to use for decoding/encoding.
// If `-p` is not provided, the default path is `area/<kind>/<kind>.hjson`.
func run() error {
	var (
		kind string
		path string
	)

	set := flag.NewFlagSet("format", flag.ContinueOnError)
	set.StringVar(&kind, "k", "", "kind of config")
	set.StringVar(&path, "p", "", "path of the config")
	if err := set.Parse(os.Args[1:]); err != nil {
		return err
	}

	cfg, ok := configs[kind]
	if !ok {
		return fmt.Errorf("%s: invalid kind", kind)
	}

	if len(path) == 0 {
		path = fmt.Sprintf("area/%s/%s.hjson", kind, kind)
	}

	if err := config.Read(path, cfg); err != nil {
		return err
	}

	if err := config.Write(path, cfg); err != nil {
		return err
	}

	return nil
}

func main() {
	logger := log.NewLogger()

	if err := run(); err != nil {
		logger.Error("could not format config", "error", err)
		os.Exit(1)
	}
}
