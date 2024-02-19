package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	python3 "go.nhat.io/python/v3"
)

func TestString(t *testing.T) {
	s := python3.NewString("Hello, World!")

	assert.NotNil(t, s)
	assert.True(t, python3.IsString(s))
	assert.Equal(t, "Hello, World!", python3.AsString(s))

	assert.False(t, python3.IsString(python3.NewList(10)))
	assert.False(t, python3.IsString(python3.NewTuple(10)))
	assert.False(t, python3.IsString(python3.NewBool(true)))
	assert.False(t, python3.IsString(python3.NewInt(42)))
	assert.False(t, python3.IsString(python3.NewInt64(42)))
	assert.True(t, python3.IsString(python3.NewString("hello")))
	assert.False(t, python3.IsString(python3.NewFloat64(3.14)))
}
