package main

import (
	"flag"
	"log"
	"os"

	"github.com/alexfalkowski/infraops/v2/internal/app/version"
)

func main() {
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
		log.Fatal(err)
	}

	if err := version.Update(name, ver, path); err != nil {
		log.Fatal(err)
	}
}
