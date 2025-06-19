package helpers

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDecodeBase64ToString(t *testing.T) {
	ctx := context.Background()

	t.Run("Successfully decode valid base64 string", func(t *testing.T) {
		// "Hello, World!" in base64
		encoded := "SGVsbG8sIFdvcmxkIQ=="
		expected := "Hello, World!"

		result := DecodeBase64ToString(ctx, encoded)

		assert.Equal(t, types.StringValue(expected), result)
	})

	t.Run("Return original string when decoding fails", func(t *testing.T) {
		// Invalid base64 string
		encoded := "This is not valid base64!"

		result := DecodeBase64ToString(ctx, encoded)

		// Should return the original string when decoding fails
		assert.Equal(t, types.StringValue(encoded), result)
	})

	t.Run("Handle empty string", func(t *testing.T) {
		encoded := ""
		expected := ""

		result := DecodeBase64ToString(ctx, encoded)

		assert.Equal(t, types.StringValue(expected), result)
	})
}

func TestByteStringToBase64(t *testing.T) {
	t.Run("Convert byte slice to base64", func(t *testing.T) {
		data := []byte("Hello, World!")
		expected := base64.StdEncoding.EncodeToString(data)

		result := ByteStringToBase64(data)

		assert.Equal(t, expected, result)
	})

	t.Run("Handle nil byte slice", func(t *testing.T) {
		var data []byte = nil

		result := ByteStringToBase64(data)

		assert.Equal(t, "", result)
	})

	t.Run("Handle empty byte slice", func(t *testing.T) {
		data := []byte{}
		expected := base64.StdEncoding.EncodeToString(data)

		result := ByteStringToBase64(data)

		assert.Equal(t, expected, result)
	})
}

func TestStringToBase64(t *testing.T) {
	t.Run("Encode string to base64", func(t *testing.T) {
		input := "Hello, World!"
		expected := base64.StdEncoding.EncodeToString([]byte(input))

		result, err := StringToBase64(input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Return error for empty string", func(t *testing.T) {
		input := ""

		result, err := StringToBase64(input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input string is empty")
		assert.Equal(t, "", result)
	})

	t.Run("Encode special characters", func(t *testing.T) {
		input := "Special chars: !@#$%^&*()"
		expected := base64.StdEncoding.EncodeToString([]byte(input))

		result, err := StringToBase64(input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Encode Unicode characters", func(t *testing.T) {
		input := "Unicode: 你好, 世界!"
		expected := base64.StdEncoding.EncodeToString([]byte(input))

		result, err := StringToBase64(input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
