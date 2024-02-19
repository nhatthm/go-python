package python

import (
	"go.nhat.io/cpy3"
)

// ErrPythonInterpreterNotInitialized is the error message when the python interpreter is not initialized.
var ErrPythonInterpreterNotInitialized = "cannot initialize the python interpreter"

var finializers = make([]func(), 0)

func init() { // nolint: gochecknoinits
	cpy3.Py_Initialize()

	if !cpy3.Py_IsInitialized() {
		panic(ErrPythonInterpreterNotInitialized)
	}
}

// Finalize finializes the python interpreter.
func Finalize() {
	for _, f := range finializers {
		f()
	}

	cpy3.Py_Finalize()
}

func registerFinalizer(f func()) {
	finializers = append(finializers, f)
}
