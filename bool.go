package python

import cpy3 "go.nhat.io/cpy/v3"

// True is a wrapper of cpy3.Py_True.
var True = NewObject(cpy3.Py_True)

// False is a wrapper of cpy3.Py_False.
var False = NewObject(cpy3.Py_False)

// IsBool returns true if the object is a Python bool.
func IsBool(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.Bool
}

// NewBool creates a new Python bool object.
func NewBool(v bool) *Object {
	return NewObject(newPyBool(v))
}

// newPyBool creates a new cpy3 Python bool object.
func newPyBool(v bool) *cpy3.PyObject {
	if v {
		return cpy3.Py_True
	}

	return cpy3.Py_False
}

// AsBool converts a Python object to a bool.
func AsBool(o *Object) bool {
	return asBool(o.PyObject())
}

func asBool(o *cpy3.PyObject) bool {
	return o.Type() == cpy3.Bool && cpy3.Py_True.RichCompareBool(o, cpy3.Py_EQ) == 1
}
