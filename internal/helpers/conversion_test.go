package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtrToString(t *testing.T) {
	t.Run("Non-nil string pointer", func(t *testing.T) {
		input := "test string"
		result := StringPtrToString(&input)
		assert.Equal(t, input, result, "Should return the dereferenced string value")
	})

	t.Run("Nil string pointer", func(t *testing.T) {
		var input *string
		result := StringPtrToString(input)
		assert.Equal(t, "", result, "Should return an empty string for nil input")
	})

	t.Run("Empty string pointer", func(t *testing.T) {
		input := ""
		result := StringPtrToString(&input)
		assert.Equal(t, "", result, "Should return an empty string for empty string input")
	})

	t.Run("String pointer with whitespace", func(t *testing.T) {
		input := "   "
		result := StringPtrToString(&input)
		assert.Equal(t, "   ", result, "Should preserve whitespace")
	})
}
