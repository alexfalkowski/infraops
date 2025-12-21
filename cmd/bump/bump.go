package main

import (
	"flag"
	"os"

	"github.com/alexfalkowski/infraops/v2/internal/app/version"
	"github.com/alexfalkowski/infraops/v2/internal/log"
)

func run() error {
	var (
		name string
		ver  string
		path string
	)

	set := flag.NewFlagSet("bump", flag.ContinueOnError)
	set.StringVar(&name, "n", "", "name of the app")
	set.StringVar(&ver, "v", "", "version of the app")
	set.StringVar(&path, "p", "", "path of the config")
	if err := set.Parse(os.Args[1:]); err != nil {
		return err
	}

	if err := version.Update(name, ver, path); err != nil {
		return err
	}

	return nil
}

func main() {
	logger := log.NewLogger()

	if err := run(); err != nil {
		logger.Error("could not bump version", "error", err)
		os.Exit(1)
	}
}
