package python3

import "go.nhat.io/cpy3"

// IsString returns true if o is a Python string object.
func IsString(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.Unicode
}

// NewString creates a new Python string object.
func NewString(s string) *Object {
	return NewObject(newPyString(s))
}

// newPyString creates a new cpy3 Python string object.
func newPyString(s string) *cpy3.PyObject {
	return cpy3.PyUnicode_FromString(s)
}

// Str returns a string representation of object o.
func Str(o *Object) string {
	return asString(o.PyObject())
}

// AsString is an alias for Str.
func AsString(o *Object) string {
	return Str(o)
}

func asString(o *cpy3.PyObject) string {
	str := o.Str()
	defer str.DecRef()

	return cpy3.PyUnicode_AsUTF8(str)
}
