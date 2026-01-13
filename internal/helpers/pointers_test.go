package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPtr(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		result := BoolPtr(true)
		assert.NotNil(t, result, "Result should not be nil")
		assert.True(t, *result, "Dereferenced value should be true")
	})

	t.Run("false value", func(t *testing.T) {
		result := BoolPtr(false)
		assert.NotNil(t, result, "Result should not be nil")
		assert.False(t, *result, "Dereferenced value should be false")
	})
}

func TestStringPtr(t *testing.T) {
	t.Run("non-empty string", func(t *testing.T) {
		input := "test string"
		result := StringPtr(input)
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, input, *result, "Dereferenced value should match input")
	})

	t.Run("empty string", func(t *testing.T) {
		result := StringPtr("")
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, "", *result, "Dereferenced value should be empty string")
	})
}

func TestGetStringValue(t *testing.T) {
	t.Run("non-nil pointer with non-empty string", func(t *testing.T) {
		input := "test value"
		ptr := &input
		result := GetStringValue(ptr)
		assert.Equal(t, input, result, "Should return the dereferenced value")
	})

	t.Run("non-nil pointer with empty string", func(t *testing.T) {
		input := ""
		ptr := &input
		result := GetStringValue(ptr)
		assert.Equal(t, "", result, "Should return empty string")
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *string = nil
		result := GetStringValue(ptr)
		assert.Equal(t, "", result, "Should return empty string for nil pointer")
	})

	t.Run("pointer to string with special characters", func(t *testing.T) {
		input := "test@#$%^&*()"
		ptr := &input
		result := GetStringValue(ptr)
		assert.Equal(t, input, result, "Should handle special characters correctly")
	})

	t.Run("pointer to multiline string", func(t *testing.T) {
		input := "line1\nline2\nline3"
		ptr := &input
		result := GetStringValue(ptr)
		assert.Equal(t, input, result, "Should handle multiline strings correctly")
	})
}
