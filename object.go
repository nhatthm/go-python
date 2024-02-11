package python3

import "go.nhat.io/cpy3"

// Objector is an interface for types that can return an Object.
type Objector interface {
	AsObject() *Object
}

// PyObjector is an interface for types that can return a PyObject.
type PyObjector interface {
	PyObject() *cpy3.PyObject
}

// Object is a wrapper around the C type python3.PyObject.
type Object cpy3.PyObject

// DecRef decreases the reference count of the object.
func (o *Object) DecRef() {
	if o == nil {
		return
	}

	(*cpy3.PyObject)(o).DecRef()
}

// PyObject returns the underlying PyObject.
func (o *Object) PyObject() *cpy3.PyObject {
	if o == nil {
		return nil
	}

	return (*cpy3.PyObject)(o)
}

// Type returns the type of the object.
func (o *Object) Type() *Object {
	return NewObject((*cpy3.PyObject)(o).Type())
}

// Length returns the length of the object.
func (o *Object) Length() int {
	return (*cpy3.PyObject)(o).Length()
}

// CallMethodArgs calls a method of the object.
func (o *Object) CallMethodArgs(name string, args ...any) *Object {
	oArgs := make([]*cpy3.PyObject, len(args))
	for i, arg := range args {
		oArgs[i] = toPyObject(arg)
	}

	return NewObject((*cpy3.PyObject)(o).CallMethodArgs(name, oArgs...))
}

// GetItem returns the item of the object.
func (o *Object) GetItem(key any) *Object {
	return NewObject((*cpy3.PyObject)(o).GetItem(toPyObject(key)))
}

// SetItem returns the item of the object.
func (o *Object) SetItem(key, value any) {
	(*cpy3.PyObject)(o).SetItem(toPyObject(key), toPyObject(value))
}

// HasItem returns true if the object has the item.
func (o *Object) HasItem(value any) bool {
	return cpy3.PySequence_Contains((*cpy3.PyObject)(o), toPyObject(value)) == 1
}

// GetAttr returns the attribute value of the object.
func (o *Object) GetAttr(name string) *Object {
	return NewObject((*cpy3.PyObject)(o).GetAttrString(name))
}

// SetAttr sets the attribute value of the object.
func (o *Object) SetAttr(name string, value any) {
	(*cpy3.PyObject)(o).SetAttrString(name, toPyObject(value))
}

// Equal returns true if the object is equal to o2.
func (o *Object) Equal(o2 *Object) bool {
	return (*cpy3.PyObject)(o).RichCompareBool((*cpy3.PyObject)(o2), cpy3.Py_EQ) == 1
}

// String returns the string representation of the object.
func (o *Object) String() string {
	return Str(o)
}

// NewObject wraps a python object with convenient methods.
func NewObject(obj *cpy3.PyObject) *Object {
	if obj == nil {
		return nil
	}

	return (*Object)(obj)
}

// toPyObject converts a value to a PyObject.
func toPyObject(v any) *cpy3.PyObject {
	return MustMarshal(v).PyObject()
}
