package python3

import (
	"reflect"

	"go.nhat.io/cpy3"
)

// IsTuple returns true if the object is a tuple.
func IsTuple(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.Tuple
}

// TupleObject is a generic Python tuple.
type TupleObject cpy3.PyObject

// DecRef decreases the reference count of the object.
func (o *TupleObject) DecRef() {
	if o == nil {
		return
	}

	(*cpy3.PyObject)(o).DecRef()
}

// PyObject returns the underlying PyObject.
func (o *TupleObject) PyObject() *cpy3.PyObject {
	return (*cpy3.PyObject)(o)
}

// Length returns the length of the tuple.
func (o *TupleObject) Length() int {
	return (*cpy3.PyObject)(o).Length()
}

// Set sets the item at index to value.
func (o *TupleObject) Set(index int, value any) {
	defer MustSuccess()

	cpy3.PyTuple_SetItem((*cpy3.PyObject)(o), index, toPyObject(value))
}

// Get returns the item at index.
func (o *TupleObject) Get(index int) *Object {
	item := cpy3.PyTuple_GetItem((*cpy3.PyObject)(o), index)

	MustSuccess()

	return NewObject(item)
}

// AsObject returns the tuple as Object.
func (o *TupleObject) AsObject() *Object {
	return (*Object)(o)
}

// String returns the string representation of the object.
func (o *TupleObject) String() string {
	return asString((*cpy3.PyObject)(o))
}

// NewTupleObject creates a new tuple.
func NewTupleObject(length int) *TupleObject {
	return (*TupleObject)(cpy3.PyTuple_New(length))
}

// Tuple is a generic Python tuple.
type Tuple[T any] struct {
	obj *TupleObject
}

// UnmarshalPyObject unmarshals a Python object to the tuple.
func (t *Tuple[T]) UnmarshalPyObject(o *Object) error {
	if !IsList(o) && !IsTuple(o) {
		return &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeOf(t)}
	}

	if IsList(o) {
		o = (*Object)((*ListObject)(o).AsTuple())
	}

	t.obj = (*TupleObject)(o)

	return nil
}

// DecRef decreases the reference count of the object.
func (t *Tuple[T]) DecRef() {
	if t == nil {
		return
	}

	t.obj.DecRef()
}

// PyObject returns the underlying PyObject.
func (t *Tuple[T]) PyObject() *cpy3.PyObject {
	return t.obj.PyObject()
}

// Length returns the length of the tuple.
func (t *Tuple[T]) Length() int {
	return t.obj.Length()
}

// Set sets the item at index to value.
func (t *Tuple[T]) Set(index int, value T) {
	t.obj.Set(index, value)
}

// Get returns the item at index.
func (t *Tuple[T]) Get(index int) T {
	o := t.obj.Get(index)

	return MustUnmarshalAs[T](o)
}

// AsObject returns the tuple as Object.
func (t *Tuple[T]) AsObject() *Object {
	return t.obj.AsObject()
}

// AsSlice converts a tuple to a slice.
func (t *Tuple[T]) AsSlice() []T {
	length := t.Length()
	slice := make([]T, length)

	for i := 0; i < length; i++ {
		slice[i] = t.Get(i)
	}

	return slice
}

// String returns the string representation of the object.
func (t *Tuple[T]) String() string {
	return t.obj.String()
}

// AnyTuple is a Python tuple.
type AnyTuple = Tuple[any]

// NewTuple creates a new tuple.
func NewTuple(length int) *AnyTuple {
	return NewTupleForType[any](length)
}

// NewTupleForType creates a new tuple for a given type.
func NewTupleForType[T any](length int) *Tuple[T] {
	return &Tuple[T]{
		obj: NewTupleObject(length),
	}
}

// NewTupleFromValues converts a slice of any data type to a tuple.
func NewTupleFromValues[T any](values ...T) *Tuple[T] {
	tuple := NewTupleForType[T](len(values))

	for i, v := range values {
		tuple.Set(i, v)
	}

	return tuple
}

// NewTupleFromAny converts a slice of any to a tuple.
func NewTupleFromAny(values ...any) *AnyTuple {
	return NewTupleFromValues(values...)
}
