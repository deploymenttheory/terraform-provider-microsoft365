package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtrToString(t *testing.T) {
	t.Run("NilPointer", func(t *testing.T) {
		var strPtr *string
		result := StringPtrToString(strPtr)
		assert.Equal(t, "", result, "Expected empty string when input is nil pointer")
	})

	t.Run("NonNilPointer", func(t *testing.T) {
		str := "test string"
		strPtr := &str
		result := StringPtrToString(strPtr)
		assert.Equal(t, "test string", result, "Expected 'test string' when input is non-nil pointer")
	})

	t.Run("EmptyStringPointer", func(t *testing.T) {
		str := ""
		strPtr := &str
		result := StringPtrToString(strPtr)
		assert.Equal(t, "", result, "Expected empty string when input is pointer to an empty string")
	})
}
