package log

import (
	"log/slog"
	"os"
)

// NewLogger returns a slog.Logger configured with a text handler writing to stdout.
//
// Callers are expected to use slog-style structured logging (key/value pairs).
func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
