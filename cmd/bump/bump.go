// Command bump updates an application's version inside the apps configuration file.
//
// By default it edits `area/apps/apps.hjson`, but you can override the location via `-p`.
//
// Usage:
//
//	bump -n <appName> -v <version> [-p <path>]
//
// Exit status is non-zero on error.
package main

import (
	"flag"
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
	set.StringVar(&name, "n", "", "name of the app")
	set.StringVar(&ver, "v", "", "version of the app")
	set.StringVar(&path, "p", "", "path of the config")
	if err := set.Parse(os.Args[1:]); err != nil {
		return err
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
		logger.Error("could not bump version", "error", err)
		os.Exit(1)
	}
}
