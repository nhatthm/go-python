package python3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/python3"
)

func TestInt(t *testing.T) {
	i := python3.NewInt(42)

	assert.NotNil(t, i)
	assert.True(t, python3.IsInt(i))
	assert.Equal(t, 42, python3.AsInt(i))

	assert.False(t, python3.IsInt(python3.NewList(10)))
	assert.False(t, python3.IsInt(python3.NewTuple(10)))
	assert.False(t, python3.IsInt(python3.NewBool(true)))
	assert.True(t, python3.IsInt(python3.NewInt(42)))
	assert.True(t, python3.IsInt(python3.NewInt64(42)))
	assert.False(t, python3.IsInt(python3.NewString("hello")))
	assert.False(t, python3.IsInt(python3.NewFloat64(3.14)))
}
