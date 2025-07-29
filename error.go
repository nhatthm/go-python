package python

import cpy3 "go.nhat.io/cpy/v3"

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
	pExec := fetchError()
	if pExec == nil {
		return nil
	}

	defer ClearError()
	defer pExec.DecRef()

	switch {
	case isException(pExec, cpy3.PyExc_ModuleNotFoundError):
		return newErrModuleNotFound(pExec)

	case isException(pExec, cpy3.PyExc_ImportError):
		return newImportError(pExec)

	case isException(pExec, cpy3.PyExc_IndexError):
		return IndexError{Exception: NewException(pExec.String())}
	}

	return NewException(pExec.String())
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
func fetchError() *Object {
	return NewObject(cpy3.PyErr_GetRaisedException())
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
