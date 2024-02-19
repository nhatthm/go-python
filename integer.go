package python

import cpy3 "go.nhat.io/cpy/v3"

// IsInt returns whether the given object is a Python int object.
func IsInt(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.Long
}

// NewInt creates a new Python int object.
func NewInt(v int) *Object {
	return NewObject(newPyInt(v))
}

// newPyInt creates a new cpy3 Python int object.
func newPyInt(v int) *cpy3.PyObject {
	return cpy3.PyLong_FromGoInt(v)
}

// NewUint creates a new Python unsigned int object.
func NewUint(v uint) *Object {
	return NewObject(newPyUint(v))
}

// newPyUint creates a new Python unsigned int object.
func newPyUint(v uint) *cpy3.PyObject {
	return cpy3.PyLong_FromGoUint(v)
}

// NewInt64 creates a new Python int object.
func NewInt64(v int64) *Object {
	return NewObject(newInt64(v))
}

// newInt64 creates a new Python int object.
func newInt64(v int64) *cpy3.PyObject {
	return cpy3.PyLong_FromGoInt64(v)
}

// NewUint64 creates a new Python unsigned int object.
func NewUint64(v uint64) *Object {
	return NewObject(newUint64(v))
}

// newUint64 creates a new Python unsigned int object.
func newUint64(v uint64) *cpy3.PyObject {
	return cpy3.PyLong_FromGoUint64(v)
}

// AsInt converts a Python object to an int.
func AsInt(o *Object) int {
	return asInt(o.PyObject())
}

func asInt(o *cpy3.PyObject) int {
	return cpy3.PyLong_AsLong(o)
}

// AsUint converts a Python object to an unsigned int.
func AsUint(o *Object) uint {
	return asUint(o.PyObject())
}

func asUint(o *cpy3.PyObject) uint {
	return cpy3.PyLong_AsUnsignedLong(o)
}

// AsInt64 converts a Python object to an int64.
func AsInt64(o *Object) int64 {
	return asInt64(o.PyObject())
}

func asInt64(o *cpy3.PyObject) int64 {
	return cpy3.PyLong_AsLongLong(o)
}

// AsUint64 converts a Python object to an unsigned int64.
func AsUint64(o *Object) uint64 {
	return asUint64(o.PyObject())
}

func asUint64(o *cpy3.PyObject) uint64 {
	return cpy3.PyLong_AsUnsignedLongLong(o)
}
