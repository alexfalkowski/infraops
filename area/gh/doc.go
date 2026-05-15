// Command gh is the Pulumi program for managing GitHub infrastructure for this repository.
//
// It reads `gh.hjson` from the current working directory (the Pulumi project directory)
// and provisions GitHub repositories and their associated configuration (for example branch
// protection, collaborators, and Pages) as described by the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/gh`, ensuring `gh.hjson` is resolved relative to that directory.
package main
