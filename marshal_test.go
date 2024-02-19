package python_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	python3 "go.nhat.io/python/v3"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		scenario       string
		value          any
		expectedResult *python3.Object
		expectedError  string
	}{
		{
			scenario:       "marshaler",
			value:          integer(42),
			expectedResult: python3.NewInt(42),
		},
		{
			scenario:       "object",
			value:          python3.NewInt(42),
			expectedResult: python3.NewInt(42),
		},
		{
			scenario:       "tuple",
			value:          python3.NewTuple(1),
			expectedResult: python3.NewTuple(1).AsObject(),
		},
		{
			scenario:       "bool",
			value:          true,
			expectedResult: python3.NewBool(true),
		},
		{
			scenario:       "string",
			value:          "hello",
			expectedResult: python3.NewString("hello"),
		},
		{
			scenario:       "rune",
			value:          'w',
			expectedResult: python3.NewInt(int('w')),
		},
		{
			scenario:       "int",
			value:          42,
			expectedResult: python3.NewInt(42),
		},
		{
			scenario:       "uint",
			value:          uint(42),
			expectedResult: python3.NewUint(42),
		},
		{
			scenario:       "int64",
			value:          int64(42),
			expectedResult: python3.NewInt64(42),
		},
		{
			scenario:       "uint64",
			value:          uint64(42),
			expectedResult: python3.NewUint64(42),
		},
		{
			scenario:       "float64",
			value:          3.14,
			expectedResult: python3.NewFloat64(3.14),
		},
		{
			scenario:       "[]int",
			value:          []int{1, 2, 3},
			expectedResult: python3.NewListFromValues(1, 2, 3).AsObject(),
		},
		{
			scenario:       "slice of custom type",
			value:          []integer{1, 2, 3},
			expectedResult: python3.NewListFromValues(integer(1), integer(2), integer(3)).AsObject(),
		},
		{
			scenario:      "unsupported",
			value:         make(chan struct{}),
			expectedError: "cannot marshal value of chan struct {} to python object",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			actual, err := python3.Marshal(tc.value)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestMustMarshal(t *testing.T) {
	assert.NotPanics(t, func() {
		python3.MustMarshal(42)
	})

	assert.Panics(t, func() {
		python3.MustMarshal(make(chan struct{}))
	})
}

func TestUnmarshal_CannotUnmarshalNonPointer(t *testing.T) {
	var actual string

	err := python3.Unmarshal(python3.NewString(""), actual)

	require.EqualError(t, err, "python3: Unmarshal(non-pointer string)")
}

func TestUnmarshal_CannotUnmarshalNil(t *testing.T) {
	err := python3.Unmarshal(python3.NewString(""), nil)

	require.EqualError(t, err, "python3: Unmarshal(nil)")
}

func TestUnmarshal_CannotUnmarshalNilInterface(t *testing.T) {
	err := python3.Unmarshal(python3.NewString(""), python3.Unmarshaler(nil))

	require.EqualError(t, err, "python3: Unmarshal(nil)")
}

func TestUnmarshal_Bool(t *testing.T) {
	testCases := []struct {
		scenario       string
		object         *python3.Object
		expectedResult bool
		expectedError  string
	}{
		{
			scenario:       "true",
			object:         python3.True,
			expectedResult: true,
		},
		{
			scenario:       "false",
			object:         python3.False,
			expectedResult: false,
		},
		{
			scenario:      "string",
			object:        python3.NewString("hello"),
			expectedError: `python3: cannot unmarshal str into Go value of type bool`,
		},
		{
			scenario:      "int",
			object:        python3.NewInt(42),
			expectedError: `python3: cannot unmarshal int into Go value of type bool`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			var actual bool

			err := python3.Unmarshal(tc.object, &actual)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestUnmarshal_String(t *testing.T) {
	testCases := []struct {
		scenario       string
		object         *python3.Object
		expectedResult string
		expectedError  string
	}{
		{
			scenario:       "empty string",
			object:         python3.NewString(""),
			expectedResult: "",
		},
		{
			scenario:       "string",
			object:         python3.NewString("hello"),
			expectedResult: "hello",
		},
		{
			scenario:      "int",
			object:        python3.NewInt(42),
			expectedError: `python3: cannot unmarshal int into Go value of type string`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			var actual string

			err := python3.Unmarshal(tc.object, &actual)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestUnmarshal_Int(t *testing.T) {
	testCases := []struct {
		scenario       string
		object         *python3.Object
		expectedResult any
		expectedError  string
	}{
		{
			scenario:       "bool",
			object:         python3.True,
			expectedResult: 0,
			expectedError:  `python3: cannot unmarshal bool into Go value of type int64`,
		},
		{
			scenario:       "empty string",
			object:         python3.NewString(""),
			expectedResult: 0,
			expectedError:  `python3: cannot unmarshal str into Go value of type int64`,
		},
		{
			scenario:       "string",
			object:         python3.NewString("hello"),
			expectedResult: 0,
			expectedError:  `python3: cannot unmarshal str into Go value of type int64`,
		},
		{
			scenario:       "float",
			object:         python3.NewFloat64(3.14),
			expectedResult: 0,
			expectedError:  `python3: cannot unmarshal float into Go value of type int64`,
		},
		{
			scenario:       "int",
			object:         python3.NewInt(42),
			expectedResult: 42,
		},
		{
			scenario:       "uint",
			object:         python3.NewUint(42),
			expectedResult: uint(42),
		},
		{
			scenario:       "int8",
			object:         python3.NewInt(42),
			expectedResult: int8(42),
		},
		{
			scenario:       "uint8",
			object:         python3.NewUint(42),
			expectedResult: uint8(42),
		},
		{
			scenario:       "int16",
			object:         python3.NewInt(42),
			expectedResult: int16(42),
		},
		{
			scenario:       "uint16",
			object:         python3.NewUint(42),
			expectedResult: uint16(42),
		},
		{
			scenario:       "int32",
			object:         python3.NewInt(42),
			expectedResult: int32(42),
		},
		{
			scenario:       "uint32",
			object:         python3.NewUint(42),
			expectedResult: uint32(42),
		},
		{
			scenario:       "int64",
			object:         python3.NewInt(42),
			expectedResult: int64(42),
		},
		{
			scenario:       "uint64",
			object:         python3.NewUint(42),
			expectedResult: uint64(42),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			actual := reflect.New(reflect.TypeOf(tc.expectedResult)).Interface()

			err := python3.Unmarshal(tc.object, actual)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestUnmarshal_Float(t *testing.T) {
	testCases := []struct {
		scenario       string
		object         *python3.Object
		expectedResult any
		expectedError  string
	}{
		{
			scenario:       "empty string",
			object:         python3.NewString(""),
			expectedResult: 0.0,
			expectedError:  `python3: cannot unmarshal str into Go value of type float64`,
		},
		{
			scenario:       "string",
			object:         python3.NewString("hello"),
			expectedResult: 0.0,
			expectedError:  `python3: cannot unmarshal str into Go value of type float64`,
		},
		{
			scenario:       "int",
			object:         python3.NewInt(42),
			expectedResult: 0.0,
			expectedError:  `python3: cannot unmarshal int into Go value of type float64`,
		},
		{
			scenario:       "float32",
			object:         python3.NewFloat64(3.14),
			expectedResult: float32(3.14),
		},
		{
			scenario:       "float64",
			object:         python3.NewFloat64(3.14),
			expectedResult: 3.14,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			actual := reflect.New(reflect.TypeOf(tc.expectedResult)).Interface()

			err := python3.Unmarshal(tc.object, actual)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestUnmarshal_Slice(t *testing.T) {
	testCases := []struct {
		scenario       string
		object         *python3.Object
		expectedResult any
		expectedError  string
	}{
		{
			scenario:       "bool",
			object:         python3.True,
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal bool into Go value of type []int`,
		},
		{
			scenario:       "string",
			object:         python3.NewString("hello"),
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal str into Go value of type []int`,
		},
		{
			scenario:       "int",
			object:         python3.NewInt(42),
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal int into Go value of type []int`,
		},
		{
			scenario:       "float",
			object:         python3.NewFloat64(3.14),
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal float into Go value of type []int`,
		},
		{
			scenario:       "empty list",
			object:         python3.NewList(0).AsObject(),
			expectedResult: []int{},
		},
		{
			scenario:       "list of string",
			object:         python3.NewListFromValues("hello", "world").AsObject(),
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal str into Go value of type int64`,
		},
		{
			scenario:       "list of int",
			object:         python3.NewListFromValues(1, 2, 3).AsObject(),
			expectedResult: []int{1, 2, 3},
		},
		{
			scenario:       "list of unmarshaler",
			object:         python3.NewListFromValues(integer(1), integer(2), integer(3)).AsObject(),
			expectedResult: []integer{1, 2, 3},
		},
		{
			scenario:       "empty tuple",
			object:         python3.NewTuple(0).AsObject(),
			expectedResult: []int{},
		},
		{
			scenario:       "tuple of string",
			object:         python3.NewTupleFromValues("hello", "world").AsObject(),
			expectedResult: []int(nil),
			expectedError:  `python3: cannot unmarshal str into Go value of type int64`,
		},
		{
			scenario:       "tuple of int",
			object:         python3.NewTupleFromValues(1, 2, 3).AsObject(),
			expectedResult: []int{1, 2, 3},
		},
		{
			scenario:       "tuple of unmarshaler",
			object:         python3.NewTupleFromValues(integer(1), integer(2), integer(3)).AsObject(),
			expectedResult: []integer{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			actual := reflect.New(reflect.TypeOf(tc.expectedResult)).Interface()

			err := python3.Unmarshal(tc.object, actual)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

type integer int

func (i integer) MarshalPyObject() *python3.Object {
	return python3.NewInt(int(i))
}

func (i *integer) UnmarshalPyObject(o *python3.Object) error { //nolint: unparam
	*i = integer(python3.AsInt(o))

	return nil
}
