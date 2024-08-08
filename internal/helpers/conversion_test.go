package helpers

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func TestTimeToString(t *testing.T) {
	t.Run("Nil time pointer", func(t *testing.T) {
		var input *time.Time
		result := TimeToString(input)
		assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")
	})

	t.Run("Valid time pointer", func(t *testing.T) {
		input := time.Date(2023, 8, 8, 12, 0, 0, 0, time.UTC)
		expected := types.StringValue(input.Format(time.RFC3339))
		result := TimeToString(&input)
		assert.Equal(t, expected, result, "Should return the time formatted as RFC3339")
	})

	t.Run("Time with different location", func(t *testing.T) {
		loc, _ := time.LoadLocation("America/New_York")
		input := time.Date(2023, 8, 8, 12, 0, 0, 0, loc)
		expected := types.StringValue(input.Format(time.RFC3339))
		result := TimeToString(&input)
		assert.Equal(t, expected, result, "Should return the time formatted as RFC3339 with correct timezone")
	})
}

func TestSliceToTypeStringSlice(t *testing.T) {
	t.Run("Nil input slice", func(t *testing.T) {
		var input []string
		result := SliceToTypeStringSlice(input)
		assert.Nil(t, result, "Should return nil for nil input")
	})

	t.Run("Empty input slice", func(t *testing.T) {
		input := []string{}
		result := SliceToTypeStringSlice(input)
		assert.Equal(t, 0, len(result), "Should return an empty slice")
	})

	t.Run("Non-empty input slice", func(t *testing.T) {
		input := []string{"one", "two", "three"}
		expected := []types.String{
			types.StringValue("one"),
			types.StringValue("two"),
			types.StringValue("three"),
		}
		result := SliceToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should convert slice of strings to slice of types.String")
	})

	t.Run("Input slice with empty strings", func(t *testing.T) {
		input := []string{"", "two", ""}
		expected := []types.String{
			types.StringValue(""),
			types.StringValue("two"),
			types.StringValue(""),
		}
		result := SliceToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should correctly handle empty strings in the input slice")
	})
}

type mockEnum string

func (e mockEnum) String() string {
	return string(e)
}

func TestEnumSliceToTypeStringSlice(t *testing.T) {
	t.Run("Nil input slice", func(t *testing.T) {
		var input []mockEnum
		result := EnumSliceToTypeStringSlice(input)
		assert.Nil(t, result, "Should return nil for nil input")
	})

	t.Run("Empty input slice", func(t *testing.T) {
		input := []mockEnum{}
		result := EnumSliceToTypeStringSlice(input)
		assert.Equal(t, 0, len(result), "Should return an empty slice")
	})

	t.Run("Non-empty input slice", func(t *testing.T) {
		input := []mockEnum{"one", "two", "three"}
		expected := []types.String{
			types.StringValue("one"),
			types.StringValue("two"),
			types.StringValue("three"),
		}
		result := EnumSliceToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should convert slice of enums to slice of types.String")
	})

	t.Run("Input slice with empty enums", func(t *testing.T) {
		input := []mockEnum{"", "two", ""}
		expected := []types.String{
			types.StringValue(""),
			types.StringValue("two"),
			types.StringValue(""),
		}
		result := EnumSliceToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should correctly handle empty enums in the input slice")
	})
}
