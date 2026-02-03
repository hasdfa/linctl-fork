package cmd

import "strings"

// contains is a test helper that wraps strings.Contains for convenience.
// It's defined in a shared file to avoid duplication across test files.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
