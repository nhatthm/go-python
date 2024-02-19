package python3_test

import (
	"os"
	"testing"

	"go.nhat.io/python3"
)

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	defer python3.Finalize()

	os.Exit(m.Run()) // nolint: gocritic
}
