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

func TestSplitCommaSeparatedString(t *testing.T) {
	t.Run("split simple comma-separated string", func(t *testing.T) {
		input := "apple,banana,cherry"
		expected := []string{"apple", "banana", "cherry"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should split string correctly")
	})

	t.Run("split single value", func(t *testing.T) {
		input := "apple"
		expected := []string{"apple"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should return single element slice")
	})

	t.Run("split empty string", func(t *testing.T) {
		input := ""
		expected := []string{}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should return empty slice for empty string")
	})

	t.Run("split string with spaces", func(t *testing.T) {
		input := "apple, banana, cherry"
		expected := []string{"apple", " banana", " cherry"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should preserve spaces after commas")
	})

	t.Run("split string with empty values", func(t *testing.T) {
		input := "apple,,cherry"
		expected := []string{"apple", "", "cherry"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should include empty strings between commas")
	})

	t.Run("split string with trailing comma", func(t *testing.T) {
		input := "apple,banana,"
		expected := []string{"apple", "banana", ""}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should include empty string for trailing comma")
	})

	t.Run("split string with leading comma", func(t *testing.T) {
		input := ",apple,banana"
		expected := []string{"", "apple", "banana"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should include empty string for leading comma")
	})

	t.Run("split string with special characters", func(t *testing.T) {
		input := "test@example.com,user#123,data$value"
		expected := []string{"test@example.com", "user#123", "data$value"}

		result := SplitCommaSeparatedString(input)

		assert.Equal(t, expected, result, "Should handle special characters correctly")
	})
}

func TestJoinWithSeparator(t *testing.T) {
	t.Run("join with comma separator", func(t *testing.T) {
		input := []string{"apple", "banana", "cherry"}
		separator := ","
		expected := "apple,banana,cherry"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should join strings with comma")
	})

	t.Run("join with space separator", func(t *testing.T) {
		input := []string{"hello", "world"}
		separator := " "
		expected := "hello world"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should join strings with space")
	})

	t.Run("join with custom separator", func(t *testing.T) {
		input := []string{"one", "two", "three"}
		separator := " | "
		expected := "one | two | three"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should join strings with custom separator")
	})

	t.Run("join single element", func(t *testing.T) {
		input := []string{"apple"}
		separator := ","
		expected := "apple"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should return single element without separator")
	})

	t.Run("join empty slice", func(t *testing.T) {
		input := []string{}
		separator := ","
		expected := ""

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should return empty string for empty slice")
	})

	t.Run("join with empty separator", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		separator := ""
		expected := "abc"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should concatenate without separator")
	})

	t.Run("join strings with empty values", func(t *testing.T) {
		input := []string{"apple", "", "cherry"}
		separator := ","
		expected := "apple,,cherry"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should include empty strings in result")
	})

	t.Run("join strings with special characters", func(t *testing.T) {
		input := []string{"test@example.com", "user#123", "data$value"}
		separator := ";"
		expected := "test@example.com;user#123;data$value"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should handle special characters correctly")
	})

	t.Run("join multiline strings", func(t *testing.T) {
		input := []string{"line1\nline2", "line3\nline4"}
		separator := " || "
		expected := "line1\nline2 || line3\nline4"

		result := JoinWithSeparator(input, separator)

		assert.Equal(t, expected, result, "Should handle multiline strings correctly")
	})
}
