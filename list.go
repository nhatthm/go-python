package python3

import (
	"reflect"

	"go.nhat.io/cpy3"
)

// IsList returns true if the object is a tuple.
func IsList(o PyObjector) bool {
	return o.PyObject().Type() == cpy3.List
}

// ListObject is a generic Python tuple.
type ListObject cpy3.PyObject

// DecRef decreases the reference count of the object.
func (o *ListObject) DecRef() {
	if o == nil {
		return
	}

	(*cpy3.PyObject)(o).DecRef()
}

// PyObject returns the underlying PyObject.
func (o *ListObject) PyObject() *cpy3.PyObject {
	return (*cpy3.PyObject)(o)
}

// Length returns the length of the tuple.
func (o *ListObject) Length() int {
	return (*cpy3.PyObject)(o).Length()
}

// Set sets the item at index to value.
func (o *ListObject) Set(index int, value any) {
	defer MustSuccess()

	cpy3.PyList_SetItem((*cpy3.PyObject)(o), index, toPyObject(value))
}

// Get returns the item at index.
func (o *ListObject) Get(index int) *Object {
	item := cpy3.PyList_GetItem((*cpy3.PyObject)(o), index)

	MustSuccess()

	return NewObject(item)
}

// AsObject returns the tuple as Object.
func (o *ListObject) AsObject() *Object {
	return (*Object)(o)
}

// AsTuple converts a list to a tuple.
func (o *ListObject) AsTuple() *TupleObject {
	return (*TupleObject)(cpy3.PyList_AsTuple((*cpy3.PyObject)(o)))
}

// String returns the string representation of the object.
func (o *ListObject) String() string {
	return asString((*cpy3.PyObject)(o))
}

// NewListObject creates a new tuple.
func NewListObject(length int) *ListObject {
	return (*ListObject)(cpy3.PyList_New(length))
}

// List is a generic Python tuple.
type List[T any] struct {
	obj *ListObject
}

func (l *List[T]) UnmarshalPyObject(o *Object) error {
	if !IsList(o) && !IsList(o) {
		return &UnmarshalTypeError{Value: TypeName(o), Type: reflect.TypeOf(l)}
	}

	if IsList(o) {

	}

	l.obj = (*ListObject)(o)

	return nil
}

// DecRef decreases the reference count of the object.
func (l *List[T]) DecRef() {
	if l == nil {
		return
	}

	l.obj.DecRef()
}

// PyObject returns the underlying PyObject.
func (l *List[T]) PyObject() *cpy3.PyObject {
	return l.obj.PyObject()
}

// Length returns the length of the tuple.
func (l *List[T]) Length() int {
	return l.obj.Length()
}

// Set sets the item at index to value.
func (l *List[T]) Set(index int, value T) {
	l.obj.Set(index, value)
}

// Get returns the item at index.
func (l *List[T]) Get(index int) T {
	o := l.obj.Get(index)

	var item T

	MustUnmarshal(o, &item)

	return item
}

// AsObject returns the tuple as Object.
func (l *List[T]) AsObject() *Object {
	return l.obj.AsObject()
}

// AsTuple converts a list to a tuple.
func (l *List[T]) AsTuple() *Tuple[T] {
	return &Tuple[T]{
		obj: l.obj.AsTuple(),
	}
}

// AsSlice converts a tuple to a slice.
func (l *List[T]) AsSlice() []T {
	length := l.Length()
	slice := make([]T, length)

	for i := 0; i < length; i++ {
		slice[i] = l.Get(i)
	}

	return slice
}

// String returns the string representation of the object.
func (l *List[T]) String() string {
	return l.obj.String()
}

// AnyList is a Python tuple.
type AnyList = List[any]

// NewList creates a new tuple.
func NewList(length int) *AnyList {
	return NewListForType[any](length)
}

// NewListForType creates a new tuple for a given type.
func NewListForType[T any](length int) *List[T] {
	return &List[T]{
		obj: NewListObject(length),
	}
}

// NewListFromValues converts a slice of any data type to a tuple.
func NewListFromValues[T any](values ...T) *List[T] {
	tuple := NewListForType[T](len(values))

	for i, v := range values {
		tuple.Set(i, v)
	}

	return tuple
}

// NewListFromAny converts a slice of any to a tuple.
func NewListFromAny(values ...any) *AnyList {
	return NewListFromValues(values...)
}
