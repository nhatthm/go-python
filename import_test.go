package python3_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/python3"
)

func TestMustImportModule(t *testing.T) {
	module := python3.MustImportModule("sys")

	attr := module.GetAttr("platform")
	defer attr.DecRef()

	actual := python3.AsString(attr)
	expected := runtime.GOOS

	assert.Equal(t, expected, actual)
}

func TestMustImportModule2_NotExists(t *testing.T) {
	err := python3.ModuleNotFoundError{
		ImportError: python3.ImportError{
			Exception: python3.Exception{
				Message: `No module named 'not_exists'`,
			},
			Module: "not_exists",
			Path:   "None",
		},
	}

	assert.PanicsWithValue(t, err, func() {
		python3.MustImportModule("not_exists")
	})
}

func TestImportModule_NotExists(t *testing.T) {
	module, actual := python3.ImportModule("not_exists")

	require.Nil(t, module)

	expected := python3.ModuleNotFoundError{
		ImportError: python3.ImportError{
			Exception: python3.Exception{
				Message: `No module named 'not_exists'`,
			},
			Module: "not_exists",
			Path:   "None",
		},
	}

	require.Equal(t, expected, actual)
	require.EqualError(t, actual, `No module named 'not_exists'`)
}
