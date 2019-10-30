package errutil

import "fmt"

// Wrap wraps error if it's not nil. Alternative to wrap from well-known github.com/pkg/errors, but for Go 1.13 usage
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(msg+": %w", err)
}
