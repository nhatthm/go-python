package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cpy3 "go.nhat.io/cpy/v3"

	python3 "go.nhat.io/python/v3"
)

func TestTuple_DecRefNil(t *testing.T) {
	var tuple *python3.AnyTuple

	assert.NotPanics(t, func() {
		tuple.DecRef()
	})
}

func TestTuple_IsTuple(t *testing.T) {
	assert.False(t, python3.IsTuple(python3.NewList(10)))
	assert.True(t, python3.IsTuple(python3.NewTuple(10)))
	assert.False(t, python3.IsTuple(python3.NewBool(true)))
	assert.False(t, python3.IsTuple(python3.NewInt(42)))
	assert.False(t, python3.IsTuple(python3.NewInt64(42)))
	assert.False(t, python3.IsTuple(python3.NewString("hello")))
	assert.False(t, python3.IsTuple(python3.NewFloat64(3.14)))
}

func TestTuple_Capacity(t *testing.T) {
	tuple := python3.NewTuple(10)
	defer tuple.DecRef()

	pyTuple := cpy3.PyTuple_New(10)
	defer pyTuple.DecRef()

	assert.True(t, cpy3.PyTuple_Check(tuple.PyObject()))
	assert.True(t, cpy3.PyTuple_CheckExact(tuple.PyObject()))
	assert.True(t, pyTuple.RichCompareBool(tuple.PyObject(), cpy3.Py_EQ) == 1) //nolint: testifylint

	assert.Equal(t, 10, tuple.Length())
}

func TestTuple_Set(t *testing.T) {
	tuple := python3.NewTuple(1)
	defer tuple.DecRef()

	assert.NotPanics(t, func() {
		tuple.Set(0, 1)

		assert.Equal(t, 1, cpy3.PyLong_AsLong(cpy3.PyTuple_GetItem(tuple.PyObject(), 0)))
	})

	err := python3.IndexError{Exception: python3.Exception{Message: `tuple assignment index out of range`}}

	assert.PanicsWithValue(t, err, func() {
		tuple.Set(1, 2)
	})
}

func TestTuple_Get(t *testing.T) {
	tuple := python3.NewTuple(1)
	defer tuple.DecRef()

	cpy3.PyTuple_SetItem(tuple.PyObject(), 0, cpy3.PyLong_FromLong(1))

	assert.NotPanics(t, func() {
		assert.Equal(t, int64(1), tuple.Get(0))
	})

	err := python3.IndexError{Exception: python3.Exception{Message: `tuple index out of range`}}

	assert.PanicsWithValue(t, err, func() {
		tuple.Get(1)
	})
}

func TestNewTupleForType(t *testing.T) {
	tuple := python3.NewTupleForType[int](3)
	defer tuple.DecRef()

	assert.Equal(t, 3, tuple.Length())

	tuple.Set(0, 3)
	tuple.Set(1, 2)
	tuple.Set(2, 1)

	assert.Equal(t, 2, tuple.Get(1))
	assert.Equal(t, []int{3, 2, 1}, tuple.AsSlice())
}

func TestNewTupleFromValues(t *testing.T) {
	tuple := python3.NewTupleFromValues(1, 2, 3)
	defer tuple.DecRef()

	assert.Equal(t, 3, tuple.Length())

	assert.Equal(t, 1, tuple.Get(0))
	assert.Equal(t, 2, tuple.Get(1))
	assert.Equal(t, 3, tuple.Get(2))

	assert.Equal(t, []int{1, 2, 3}, tuple.AsSlice())
}

func TestNewTupleFromAny(t *testing.T) {
	tuple := python3.NewTupleFromAny(1, "hello", 3.14)
	defer tuple.DecRef()

	assert.Equal(t, 3, tuple.Length())

	assert.Equal(t, int64(1), tuple.Get(0))
	assert.Equal(t, "hello", tuple.Get(1))
	assert.InDelta(t, 3.14, tuple.Get(2), 0.01)

	assert.Equal(t, []any{int64(1), "hello", 3.14}, tuple.AsSlice())
}

func TestTupleOfList(t *testing.T) {
	tuple := python3.NewTupleForType[[]int](1)

	tuple.Set(0, []int{1, 2})

	assert.Equal(t, `([1, 2],)`, tuple.String())

	actual := tuple.Get(0)
	expected := []int{1, 2}

	assert.Equal(t, expected, actual)
}

func TestTupleOfTuple(t *testing.T) {
	tuple := python3.NewTupleForType[*python3.Tuple[int]](1)

	tuple.Set(0, python3.NewTupleFromValues(1, 2))

	assert.Equal(t, `((1, 2),)`, tuple.String())

	actual := tuple.Get(0)
	expected := python3.NewTupleFromValues(1, 2)

	assert.True(t, expected.AsObject().Equal(actual.AsObject()))
}
