// Package agentutil holds helpers shared by cookbook example programs.
package agentutil

import (
	"os"
	"strings"
)

// EnvOr returns the trimmed environment variable or fallback when unset.
func EnvOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// IsTerminal reports whether f is a character device (TTY).
func IsTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Compact collapses whitespace in s to single spaces.
func Compact(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
