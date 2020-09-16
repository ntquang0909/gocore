package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	assert.Equal(t, BoolValue(Bool(true)), true)

	assert.Equal(t, StringValue(String("Test")), "Test")

	assert.Equal(t, StringSliceValue(StringSlice([]string{"Test", "Test"})), []string{"Test", "Test"})

	assert.Equal(t, Int64Value(Int64(99)), int64(99))

	assert.Equal(t, IntValue(Int(99)), int(99))

	assert.Equal(t, Float64Value(Float64(99)), float64(99))

	assert.Equal(t, Float32Value(Float32(99)), float32(99))
}
