package errors

import (
	"fmt"
)

// Prefix the error.
func Prefix(p string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%v: %w", p, err)
}
