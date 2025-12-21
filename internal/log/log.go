package log

import (
	"log/slog"
	"os"
)

// NewLogger returns a text logger.
func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
