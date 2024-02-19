package python3

import (
	"go.nhat.io/cpy3"
)

// Exception is a Python exception.
type Exception struct { //nolint: errname,stylecheck
	Message string
}

// Error returns a string representation of the Exception.
func (e Exception) Error() string {
	return e.Message
}

// ImportError is returned when a Python module cannot be imported.
type ImportError struct {
	Exception

	Module string
	Path   string
}

// ModuleNotFoundError is returned when a Python module cannot be found.
type ModuleNotFoundError struct {
	ImportError
}

// IndexError is returned when a sequence subscript is out of range.
type IndexError struct {
	Exception
}

// MustSuccess panics if the last Python operation failed.
func MustSuccess() {
	if err := LastError(); err != nil {
		panic(err)
	}
}

// LastError returns the last error that occurred in the Python interpreter.
func LastError() error {
	pType, pValue, pTraceback := fetchError()
	if pValue == nil {
		return nil
	}

	defer ClearError()
	defer pType.DecRef()
	defer pValue.DecRef()
	defer pTraceback.DecRef()

	switch {
	case isException(pType, cpy3.PyExc_ModuleNotFoundError):
		return newErrModuleNotFound(pValue)

	case isException(pType, cpy3.PyExc_ImportError):
		return newImportError(pValue)

	case isException(pType, cpy3.PyExc_IndexError):
		return IndexError{Exception: NewException(pValue.String())}
	}

	return NewException(pValue.String())
}

// isException returns true if err is an instance of ex.
func isException(err *Object, target *cpy3.PyObject) bool {
	if target == nil {
		return err.PyObject() == target
	}

	return cpy3.PyErr_GivenExceptionMatches(err.PyObject(), target)
}

// ClearError clears the last error that occurred in the Python interpreter.
func ClearError() {
	cpy3.PyErr_Clear()
}

// fetchError returns the last error that occurred in the Python interpreter.
func fetchError() (*Object, *Object, *Object) {
	pType, pValue, pTraceback := cpy3.PyErr_Fetch()
	if pType == nil {
		return nil, nil, nil
	}

	return NewObject(pType), NewObject(pValue), NewObject(pTraceback)
}

// NewException creates a new Exception.
func NewException(message string) Exception {
	return Exception{Message: message}
}

func newErrModuleNotFound(err *Object) ModuleNotFoundError {
	return ModuleNotFoundError{
		ImportError: newImportError(err),
	}
}

func newImportError(err *Object) ImportError {
	name := err.GetAttr("name")
	path := err.GetAttr("path")
	msg := err.GetAttr("msg")

	defer name.DecRef()
	defer path.DecRef()
	defer msg.DecRef()

	return ImportError{
		Exception: NewException(msg.String()),
		Module:    name.String(),
		Path:      path.String(),
	}
}
