package python3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/python3"
)

func TestBool(t *testing.T) {
	boolT := python3.NewBool(true)
	boolF := python3.NewBool(false)

	assert.True(t, python3.IsBool(boolT))
	assert.True(t, python3.IsBool(boolF))
	assert.False(t, python3.IsBool(python3.NewInt(0)))
	assert.False(t, python3.IsBool(python3.NewString("")))

	assert.True(t, python3.AsBool(boolT))
	assert.False(t, python3.AsBool(boolF))

	assert.False(t, python3.AsBool(python3.NewInt(1)))
	assert.False(t, python3.AsBool(python3.NewInt(0)))

	assert.False(t, python3.IsBool(python3.NewList(10)))
	assert.False(t, python3.IsBool(python3.NewTuple(10)))
	assert.False(t, python3.IsBool(python3.NewInt(42)))
	assert.False(t, python3.IsBool(python3.NewInt64(42)))
	assert.False(t, python3.IsBool(python3.NewString("hello")))
	assert.False(t, python3.IsBool(python3.NewFloat64(3.14)))
}
