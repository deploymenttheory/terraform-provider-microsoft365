package state

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

func TestBoolPtrToTypeBool(t *testing.T) {
	t.Run("Nil bool pointer", func(t *testing.T) {
		var input *bool
		result := BoolPtrToTypeBool(input)
		assert.True(t, result.IsNull(), "Should return types.BoolNull() for nil input")
	})

	t.Run("True bool pointer", func(t *testing.T) {
		input := true
		result := BoolPtrToTypeBool(&input)
		assert.Equal(t, types.BoolValue(true), result, "Should return types.BoolValue(true) for true input")
	})

	t.Run("False bool pointer", func(t *testing.T) {
		input := false
		result := BoolPtrToTypeBool(&input)
		assert.Equal(t, types.BoolValue(false), result, "Should return types.BoolValue(false) for false input")
	})
}

type testEnum int

const (
	EnumOne testEnum = iota
	EnumTwo
	EnumThree
)

func (e testEnum) String() string {
	return [...]string{"One", "Two", "Three"}[e]
}

func TestEnumPtrToTypeString(t *testing.T) {
	t.Run("Nil enum pointer", func(t *testing.T) {
		var input *testEnum
		result := EnumPtrToTypeString(input)
		assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")
	})

	t.Run("Valid enum pointer", func(t *testing.T) {
		input := EnumTwo
		result := EnumPtrToTypeString(&input)
		assert.Equal(t, types.StringValue("Two"), result, "Should return the string representation of the enum")
	})

	t.Run("Different enum values", func(t *testing.T) {
		testCases := []struct {
			input    testEnum
			expected string
		}{
			{EnumOne, "One"},
			{EnumTwo, "Two"},
			{EnumThree, "Three"},
		}

		for _, tc := range testCases {
			result := EnumPtrToTypeString(&tc.input)
			assert.Equal(t, types.StringValue(tc.expected), result, "Should return correct string for enum value")
		}
	})
}

func TestInt32PtrToTypeInt64(t *testing.T) {
	t.Run("Nil int32 pointer", func(t *testing.T) {
		var input *int32
		result := Int32PtrToTypeInt64(input)
		assert.True(t, result.IsNull(), "Should return types.Int64Null() for nil input")
	})

	t.Run("Valid int32 pointer", func(t *testing.T) {
		input := int32(42)
		result := Int32PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(42), result, "Should return types.Int64Value(42) for input 42")
	})

	t.Run("Negative int32 pointer", func(t *testing.T) {
		input := int32(-123)
		result := Int32PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(-123), result, "Should return types.Int64Value(-123) for input -123")
	})

	t.Run("Zero int32 pointer", func(t *testing.T) {
		input := int32(0)
		result := Int32PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(0), result, "Should return types.Int64Value(0) for input 0")
	})

	t.Run("Max int32 pointer", func(t *testing.T) {
		input := int32(2147483647) // Max value for int32
		result := Int32PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(2147483647), result, "Should correctly convert max int32 value")
	})

	t.Run("Min int32 pointer", func(t *testing.T) {
		input := int32(-2147483648) // Min value for int32
		result := Int32PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(-2147483648), result, "Should correctly convert min int32 value")
	})
}
