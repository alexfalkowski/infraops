package main

import (
	"errors"
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
	set.Usage = func() {
		fmt.Fprintf(set.Output(), "Usage: %s -k <apps|cf|do|gh> [-p <path>]\n\n", set.Name())
		fmt.Fprintln(set.Output(), "Normalizes an area HJSON configuration file in place.")
		fmt.Fprintln(set.Output(), "By default, the path is area/<kind>/<kind>.hjson.")
		fmt.Fprintln(set.Output(), "\nFlags:")
		set.PrintDefaults()
	}
	set.StringVar(&kind, "k", "", "config kind (apps|cf|do|gh) (required)")
	set.StringVar(&path, "p", "", "config file path (default area/<kind>/<kind>.hjson)")
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
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		logger.Error("could not format config", "error", err)
		os.Exit(1)
	}
}
