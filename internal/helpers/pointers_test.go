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
