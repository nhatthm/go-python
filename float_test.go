package python_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	python3 "go.nhat.io/python/v3"
)

func TestFloat64(t *testing.T) {
	f := python3.NewFloat64(3.14)

	assert.NotNil(t, f)
	assert.True(t, python3.IsFloat(f))
	assert.InDelta(t, 3.14, python3.AsFloat64(f), 0.01)

	assert.False(t, python3.IsFloat(python3.NewList(10)))
	assert.False(t, python3.IsFloat(python3.NewTuple(10)))
	assert.False(t, python3.IsFloat(python3.NewBool(true)))
	assert.False(t, python3.IsFloat(python3.NewInt(42)))
	assert.False(t, python3.IsFloat(python3.NewString("hello")))
}
