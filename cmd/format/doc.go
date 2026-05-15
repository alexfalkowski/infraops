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
