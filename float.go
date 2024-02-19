package python

import "go.nhat.io/cpy3"

// IsFloat returns true if the object is a Python float.
func IsFloat(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.Float
}

// NewFloat64 creates a new Python float object.
func NewFloat64(v float64) *Object {
	return NewObject(newPyFloat64(v))
}

// newPyFloat64 creates a new cpy3 Python float object.
func newPyFloat64(v float64) *cpy3.PyObject {
	return cpy3.PyFloat_FromDouble(v)
}

// AsFloat64 converts a Python object to a float64.
func AsFloat64(o *Object) float64 {
	return asFloat64(o.PyObject())
}

func asFloat64(o *cpy3.PyObject) float64 {
	return cpy3.PyFloat_AsDouble(o)
}
