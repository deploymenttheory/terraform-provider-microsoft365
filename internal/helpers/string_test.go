package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractGUIDFromString(t *testing.T) {
	t.Run("valid GUID at beginning", func(t *testing.T) {
		input := "12345678-1234-1234-1234-123456789abc"
		expected := "12345678-1234-1234-1234-123456789abc"

		result, err := ExtractGUIDFromString(input)

		assert.NoError(t, err, "Should not return an error for valid GUID")
		assert.Equal(t, expected, result, "Should extract the correct GUID")
	})

	t.Run("valid GUID with suffix", func(t *testing.T) {
		input := "12345678-1234-1234-1234-123456789abc_additional_text"
		expected := "12345678-1234-1234-1234-123456789abc"

		result, err := ExtractGUIDFromString(input)

		assert.NoError(t, err, "Should not return an error for valid GUID with suffix")
		assert.Equal(t, expected, result, "Should extract only the GUID part")
	})

	t.Run("invalid GUID format", func(t *testing.T) {
		input := "invalid-guid-format"

		result, err := ExtractGUIDFromString(input)

		assert.Error(t, err, "Should return an error for invalid GUID format")
		assert.Empty(t, result, "Result should be empty for invalid GUID")
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""

		result, err := ExtractGUIDFromString(input)

		assert.Error(t, err, "Should return an error for empty string")
		assert.Empty(t, result, "Result should be empty for empty string")
	})

	t.Run("GUID not at beginning", func(t *testing.T) {
		input := "prefix_12345678-1234-1234-1234-123456789abc"

		result, err := ExtractGUIDFromString(input)

		assert.Error(t, err, "Should return an error when GUID is not at beginning")
		assert.Empty(t, result, "Result should be empty when GUID is not at beginning")
	})
}

func TestStringToInt(t *testing.T) {
	t.Run("valid string mapping", func(t *testing.T) {
		mapping := map[string]int{
			"low":    1,
			"medium": 2,
			"high":   3,
		}

		result, err := StringToInt("medium", mapping)

		assert.NoError(t, err, "Should not return an error for valid string")
		assert.Equal(t, 2, result, "Should return the correct integer value")
	})

	t.Run("invalid string", func(t *testing.T) {
		mapping := map[string]int{
			"low":    1,
			"medium": 2,
			"high":   3,
		}

		result, err := StringToInt("critical", mapping)

		assert.Error(t, err, "Should return an error for invalid string")
		assert.Equal(t, -1, result, "Should return -1 for invalid string")
	})

	t.Run("empty string", func(t *testing.T) {
		mapping := map[string]int{
			"low":    1,
			"medium": 2,
			"high":   3,
		}

		result, err := StringToInt("", mapping)

		assert.Error(t, err, "Should return an error for empty string")
		assert.Equal(t, -1, result, "Should return -1 for empty string")
	})

	t.Run("empty mapping", func(t *testing.T) {
		mapping := map[string]int{}

		result, err := StringToInt("test", mapping)

		assert.Error(t, err, "Should return an error for empty mapping")
		assert.Equal(t, -1, result, "Should return -1 for empty mapping")
	})
}
