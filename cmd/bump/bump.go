package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/alexfalkowski/infraops/v2/internal/app/version"
	"github.com/alexfalkowski/infraops/v2/internal/log"
)

// run parses CLI flags and updates the application version in the configuration file.
func run() error {
	var (
		name string
		ver  string
		path string
	)

	set := flag.NewFlagSet("bump", flag.ContinueOnError)
	set.Usage = func() {
		fmt.Fprintf(set.Output(), "Usage: %s -n <appName> -v <version> [-p <path>]\n\n", set.Name())
		fmt.Fprintln(set.Output(), "Updates an application version in place.")
		fmt.Fprintln(set.Output(), "By default, the path is area/apps/apps.hjson.")
		fmt.Fprintln(set.Output(), "The version is written exactly as provided and is expected to be semantic.")
		fmt.Fprintln(set.Output(), "\nFlags:")
		set.PrintDefaults()
	}
	set.StringVar(&name, "n", "", "application name (required)")
	set.StringVar(&ver, "v", "", "application version (required)")
	set.StringVar(&path, "p", "", "config file path (default area/apps/apps.hjson)")
	if err := set.Parse(os.Args[1:]); err != nil {
		return err
	}

	if name == "" {
		return errors.New("application name is required")
	}

	if ver == "" {
		return errors.New("application version is required")
	}

	if len(path) == 0 {
		path = "area/apps/apps.hjson"
	}

	if err := version.Update(name, ver, path); err != nil {
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

		logger.Error("could not bump version", "error", err)
		os.Exit(1)
	}
}
