// Command cf is the Pulumi program for managing Cloudflare infrastructure for this repository.
//
// It reads `cf.hjson` from the current working directory (the Pulumi project directory) and
// provisions Cloudflare resources (for example zones, DNS records, and R2 buckets/custom domains)
// as described by the config.
//
// This program is typically executed via the repository Makefile targets, which run Pulumi
// with `--cwd area/cf`, ensuring `cf.hjson` is resolved relative to that directory.
package main
