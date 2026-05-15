// Command bump updates an application's version inside the apps configuration file.
//
// By default it edits `area/apps/apps.hjson`, but you can override the location via `-p`.
// The tool is intended for internal automation and expects `-v` to be a semantic version string
// supplied by that automation.
//
// Usage:
//
//	bump -n <appName> -v <version> [-p <path>]
//
// Exit status is non-zero on error.
package main
