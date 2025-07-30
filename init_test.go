package python_test

import (
	"os"
	"runtime"
	"testing"

	python3 "go.nhat.io/python/v3"
)

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	runtime.LockOSThread()

	ret := m.Run()

	python3.Finalize()
	runtime.UnlockOSThread()

	os.Exit(ret)
}
