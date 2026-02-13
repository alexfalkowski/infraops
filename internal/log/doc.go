// Package log provides a minimal slog logger configuration for repository CLI tools.
//
// The logger is intentionally simple: it writes human-readable text logs to stdout with
// the default handler options.
//
// This package is used by small command binaries (for example under cmd/) that need
// consistent structured logging without pulling in a heavier logging framework.
package log
