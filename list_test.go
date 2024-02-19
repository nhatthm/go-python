package python3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.nhat.io/cpy3"

	"go.nhat.io/python3"
)

func TestList_DecRefNil(t *testing.T) {
	var list *python3.AnyList

	assert.NotPanics(t, func() {
		list.DecRef()
	})
}

func TestList_IsList(t *testing.T) {
	assert.True(t, python3.IsList(python3.NewList(10)))
	assert.False(t, python3.IsList(python3.NewTuple(10)))
	assert.False(t, python3.IsList(python3.NewBool(true)))
	assert.False(t, python3.IsList(python3.NewInt(42)))
	assert.False(t, python3.IsList(python3.NewInt64(42)))
	assert.False(t, python3.IsList(python3.NewString("hello")))
	assert.False(t, python3.IsList(python3.NewFloat64(3.14)))
}

func TestList_Capacity(t *testing.T) {
	list := python3.NewList(10)
	defer list.DecRef()

	pyList := cpy3.PyList_New(10)
	defer pyList.DecRef()

	assert.True(t, cpy3.PyList_Check(list.PyObject()))
	assert.True(t, cpy3.PyList_CheckExact(list.PyObject()))
	assert.True(t, pyList.RichCompareBool(list.PyObject(), cpy3.Py_EQ) == 1) //nolint: testifylint

	assert.Equal(t, 10, list.Length())
}

func TestList_Set(t *testing.T) {
	list := python3.NewList(1)
	defer list.DecRef()

	assert.NotPanics(t, func() {
		list.Set(0, 1)

		assert.Equal(t, 1, cpy3.PyLong_AsLong(cpy3.PyList_GetItem(list.PyObject(), 0)))
	})

	err := python3.IndexError{Exception: python3.Exception{Message: `list assignment index out of range`}}

	assert.PanicsWithValue(t, err, func() {
		list.Set(1, 2)
	})
}

func TestList_Get(t *testing.T) {
	list := python3.NewList(1)
	defer list.DecRef()

	cpy3.PyList_SetItem(list.PyObject(), 0, cpy3.PyLong_FromLong(1))

	assert.NotPanics(t, func() {
		assert.Equal(t, int64(1), list.Get(0))
	})

	err := python3.IndexError{Exception: python3.Exception{Message: `list index out of range`}}

	assert.PanicsWithValue(t, err, func() {
		list.Get(1)
	})
}

func TestNewListForType(t *testing.T) {
	list := python3.NewListForType[int](3)
	defer list.DecRef()

	assert.Equal(t, 3, list.Length())

	list.Set(0, 3)
	list.Set(1, 2)
	list.Set(2, 1)

	assert.Equal(t, 2, list.Get(1))
	assert.Equal(t, []int{3, 2, 1}, list.AsSlice())
}

func TestNewListFromValues(t *testing.T) {
	list := python3.NewListFromValues(1, 2, 3)
	defer list.DecRef()

	assert.Equal(t, 3, list.Length())

	assert.Equal(t, 1, list.Get(0))
	assert.Equal(t, 2, list.Get(1))
	assert.Equal(t, 3, list.Get(2))

	assert.Equal(t, []int{1, 2, 3}, list.AsSlice())
}

func TestNewListFromAny(t *testing.T) {
	list := python3.NewListFromAny(1, "hello", 3.14)
	defer list.DecRef()

	assert.Equal(t, 3, list.Length())

	assert.Equal(t, int64(1), list.Get(0))
	assert.Equal(t, "hello", list.Get(1))
	assert.InDelta(t, 3.14, list.Get(2), 0.01)

	assert.Equal(t, []any{int64(1), "hello", 3.14}, list.AsSlice())
}

func TestListOfListOfList(t *testing.T) {
	list := python3.NewListForType[[][]int](1)

	list.Set(0, [][]int{{1, 2}, {2, 3}})

	actual := list.Get(0)
	expected := [][]int{{1, 2}, {2, 3}}

	assert.Equal(t, expected, actual)
}

func TestListOfTuple(t *testing.T) {
	list := python3.NewListForType[*python3.Tuple[int]](1)

	list.Set(0, python3.NewTupleFromValues(1, 2))

	assert.Equal(t, `[(1, 2)]`, list.AsObject().String())

	actual := list.Get(0)
	expected := python3.NewTupleFromValues(1, 2)

	assert.True(t, expected.AsObject().Equal(actual.AsObject()))
}

func TestListFromTuple(t *testing.T) {
	tuple := python3.NewTupleFromValues(1, 2, 3)

	var list python3.List[int]

	err := list.UnmarshalPyObject(tuple.AsObject())
	require.NoError(t, err)

	assert.Equal(t, `[1, 2, 3]`, list.String())
}
