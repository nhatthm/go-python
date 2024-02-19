package python_test

import (
	"os"
	"testing"

	python3 "go.nhat.io/python/v3"
)

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	defer python3.Finalize()

	os.Exit(m.Run()) // nolint: gocritic
}
