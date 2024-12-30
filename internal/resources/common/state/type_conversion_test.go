package state

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
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

type enumOneType int

const (
	EnumOne enumOneType = iota
	EnumTwo
	EnumThree
)

func (e enumOneType) String() string {
	return [...]string{"One", "Two", "Three"}[e]
}

func TestEnumPtrToTypeString(t *testing.T) {
	t.Run("Nil enum pointer", func(t *testing.T) {
		var input *enumOneType
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
			input    enumOneType
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

type enumTwoType int

const (
	EnumApple enumTwoType = iota
	EnumBanana
	EnumCherry
)

func (e enumTwoType) String() string {
	return [...]string{"Apple", "Banana", "Cherry"}[e]
}

func TestEnumListPtrToTypeStringSlice(t *testing.T) {
	t.Run("Nil input slice", func(t *testing.T) {
		var input []*enumTwoType
		result := EnumListPtrToTypeStringSlice(input)
		assert.Nil(t, result, "Should return nil for nil input")
	})

	t.Run("Empty input slice", func(t *testing.T) {
		input := []*enumTwoType{}
		result := EnumListPtrToTypeStringSlice(input)
		assert.Equal(t, 0, len(result), "Should return an empty slice")
	})

	t.Run("Non-empty input slice with valid enum pointers", func(t *testing.T) {
		apple := EnumApple
		banana := EnumBanana
		input := []*enumTwoType{&apple, &banana}
		expected := []types.String{
			types.StringValue("Apple"),
			types.StringValue("Banana"),
		}
		result := EnumListPtrToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should convert slice of enum pointers to slice of types.String")
	})

	t.Run("Input slice with nil enum pointers", func(t *testing.T) {
		apple := EnumApple
		input := []*enumTwoType{&apple, nil}
		expected := []types.String{
			types.StringValue("Apple"),
			types.StringNull(),
		}
		result := EnumListPtrToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should handle nil pointers correctly")
	})

	t.Run("Input slice with all nil enum pointers", func(t *testing.T) {
		input := []*enumTwoType{nil, nil}
		expected := []types.String{
			types.StringNull(),
			types.StringNull(),
		}
		result := EnumListPtrToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should handle slice with only nil pointers correctly")
	})

	t.Run("Input slice with different enum values", func(t *testing.T) {
		apple := EnumApple
		cherry := EnumCherry
		input := []*enumTwoType{&apple, &cherry}
		expected := []types.String{
			types.StringValue("Apple"),
			types.StringValue("Cherry"),
		}
		result := EnumListPtrToTypeStringSlice(input)
		assert.Equal(t, expected, result, "Should convert multiple different enum values correctly")
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

func TestInt32PtrToTypeInt32(t *testing.T) {
	t.Run("Nil int32 pointer", func(t *testing.T) {
		var input *int32
		result := Int32PtrToTypeInt32(input)
		assert.True(t, result.IsNull(), "Should return types.Int32Null() for nil input")
	})

	t.Run("Valid int32 pointer", func(t *testing.T) {
		input := int32(42)
		result := Int32PtrToTypeInt32(&input)
		assert.Equal(t, types.Int32Value(42), result, "Should return types.Int32Value(42) for input 42")
	})

	t.Run("Negative int32 pointer", func(t *testing.T) {
		input := int32(-123)
		result := Int32PtrToTypeInt32(&input)
		assert.Equal(t, types.Int32Value(-123), result, "Should return types.Int32Value(-123) for input -123")
	})

	t.Run("Zero int32 pointer", func(t *testing.T) {
		input := int32(0)
		result := Int32PtrToTypeInt32(&input)
		assert.Equal(t, types.Int32Value(0), result, "Should return types.Int32Value(0) for input 0")
	})

	t.Run("Max int32 pointer", func(t *testing.T) {
		input := int32(2147483647) // Max value for int32
		result := Int32PtrToTypeInt32(&input)
		assert.Equal(t, types.Int32Value(2147483647), result, "Should correctly convert max int32 value")
	})

	t.Run("Min int32 pointer", func(t *testing.T) {
		input := int32(-2147483648) // Min value for int32
		result := Int32PtrToTypeInt32(&input)
		assert.Equal(t, types.Int32Value(-2147483648), result, "Should correctly convert min int32 value")
	})
}

func TestDateOnlyPtrToString(t *testing.T) {
	t.Run("Nil DateOnly pointer", func(t *testing.T) {
		var input *serialization.DateOnly
		result := DateOnlyPtrToString(input)
		assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")
	})

	t.Run("Valid DateOnly pointer", func(t *testing.T) {
		date := time.Date(2024, 8, 16, 0, 0, 0, 0, time.UTC)
		input := serialization.NewDateOnly(date)
		expected := types.StringValue("2024-08-16")
		result := DateOnlyPtrToString(input)
		assert.Equal(t, expected, result, "Should return the date formatted as YYYY-MM-DD")
	})

	t.Run("Different valid DateOnly pointer", func(t *testing.T) {
		date := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
		input := serialization.NewDateOnly(date)
		expected := types.StringValue("1999-12-31")
		result := DateOnlyPtrToString(input)
		assert.Equal(t, expected, result, "Should return the date formatted as YYYY-MM-DD")
	})

	t.Run("Edge case DateOnly pointer", func(t *testing.T) {
		date := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) // The minimum date in the Gregorian calendar
		input := serialization.NewDateOnly(date)
		expected := types.StringValue("0001-01-01")
		result := DateOnlyPtrToString(input)
		assert.Equal(t, expected, result, "Should handle the minimum date correctly")
	})

	t.Run("Another edge case DateOnly pointer", func(t *testing.T) {
		date := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC) // The maximum date in the Gregorian calendar
		input := serialization.NewDateOnly(date)
		expected := types.StringValue("9999-12-31")
		result := DateOnlyPtrToString(input)
		assert.Equal(t, expected, result, "Should handle the maximum date correctly")
	})
}

func TestByteToString(t *testing.T) {
	t.Run("Empty byte slice", func(t *testing.T) {
		input := []byte{}
		result := ByteToString(input)
		assert.Equal(t, "", result, "Should return an empty string for empty byte slice")
	})

	t.Run("Non-empty byte slice", func(t *testing.T) {
		input := []byte("Hello, World!")
		expected := "SGVsbG8sIFdvcmxkIQ==" // Base64 encoded "Hello, World!"
		result := ByteToString(input)
		assert.Equal(t, expected, result, "Should return base64 encoded string")
	})

	t.Run("Byte slice with special characters", func(t *testing.T) {
		input := []byte("Hello, 世界!")
		expected := "SGVsbG8sIOS4lueVjCE=" // Base64 encoded "Hello, 世界!"
		result := ByteToString(input)
		assert.Equal(t, expected, result, "Should correctly encode special characters")
	})

	t.Run("Byte slice with null bytes", func(t *testing.T) {
		input := []byte{0, 1, 2, 3}
		expected := "AAECAw==" // Base64 encoded [0, 1, 2, 3]
		result := ByteToString(input)
		assert.Equal(t, expected, result, "Should correctly encode null bytes")
	})
}

func TestBoolPtrToBool(t *testing.T) {
	t.Run("Nil bool pointer", func(t *testing.T) {
		var input *bool
		result := BoolPtrToBool(input)
		assert.False(t, result, "Should return false for nil input")
	})

	t.Run("True bool pointer", func(t *testing.T) {
		input := true
		result := BoolPtrToBool(&input)
		assert.True(t, result, "Should return true for true input")
	})

	t.Run("False bool pointer", func(t *testing.T) {
		input := false
		result := BoolPtrToBool(&input)
		assert.False(t, result, "Should return false for false input")
	})
}

func TestInt64PtrToTypeInt64(t *testing.T) {
	t.Run("Nil int64 pointer", func(t *testing.T) {
		var input *int64
		result := Int64PtrToTypeInt64(input)
		assert.True(t, result.IsNull(), "Should return types.Int64Null() for nil input")
	})

	t.Run("Valid int64 pointer", func(t *testing.T) {
		input := int64(42)
		result := Int64PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(42), result, "Should return types.Int64Value(42) for input 42")
	})

	t.Run("Negative int64 pointer", func(t *testing.T) {
		input := int64(-123)
		result := Int64PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(-123), result, "Should return types.Int64Value(-123) for input -123")
	})

	t.Run("Zero int64 pointer", func(t *testing.T) {
		input := int64(0)
		result := Int64PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(0), result, "Should return types.Int64Value(0) for input 0")
	})

	t.Run("Max int64 pointer", func(t *testing.T) {
		input := int64(9223372036854775807) // Max value for int64
		result := Int64PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(9223372036854775807), result, "Should correctly convert max int64 value")
	})

	t.Run("Min int64 pointer", func(t *testing.T) {
		input := int64(-9223372036854775808) // Min value for int64
		result := Int64PtrToTypeInt64(&input)
		assert.Equal(t, types.Int64Value(-9223372036854775808), result, "Should correctly convert min int64 value")
	})
}

func TestISO8601DurationToString(t *testing.T) {
	t.Run("Nil ISODuration pointer", func(t *testing.T) {
		var input *serialization.ISODuration
		result := ISO8601DurationToString(input)
		assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")
	})

	t.Run("Valid ISODuration pointer with years", func(t *testing.T) {
		input := serialization.NewDuration(1, 0, 3, 4, 5, 6, 7) // Example duration: P1Y3DT4H5M6S
		expected := "P1Y3DT4H5M6S"
		result := ISO8601DurationToString(input)
		assert.Equal(t, types.StringValue(expected), result, "Should correctly convert valid ISODuration to ISO 8601 string")
	})

	t.Run("Valid ISODuration pointer with weeks", func(t *testing.T) {
		input := serialization.NewDuration(0, 2, 0, 0, 0, 0, 0) // Example duration: P2W
		expected := "P2W"
		result := ISO8601DurationToString(input)
		assert.Equal(t, types.StringValue(expected), result, "Should correctly convert ISODuration with weeks to ISO 8601 string")
	})

	t.Run("Empty ISODuration", func(t *testing.T) {
		input := serialization.NewDuration(0, 0, 0, 0, 0, 0, 0) // Example duration: P
		expected := "P"
		result := ISO8601DurationToString(input)
		assert.Equal(t, types.StringValue(expected), result, "Should correctly convert an empty ISODuration to 'P'")
	})
}

func TestDecodeBase64ToString(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid base64 string", func(t *testing.T) {
		input := base64.StdEncoding.EncodeToString([]byte("test content")) // Encodes "test content" to base64
		expected := "test content"
		result := DecodeBase64ToString(ctx, input)
		assert.Equal(t, types.StringValue(expected), result, "Should return the decoded content as a types.String")
	})

	t.Run("Invalid base64 string", func(t *testing.T) {
		input := "invalid_base64" // Not a valid base64 string
		result := DecodeBase64ToString(ctx, input)
		assert.Equal(t, types.StringValue(input), result, "Should return the original string as a types.String on decoding failure")
	})

	t.Run("Empty base64 string", func(t *testing.T) {
		input := ""    // Empty string
		expected := "" // Decoded result of an empty base64 string is also an empty string
		result := DecodeBase64ToString(ctx, input)
		assert.Equal(t, types.StringValue(expected), result, "Should return an empty string as a types.String")
	})

	t.Run("Valid base64 with padding", func(t *testing.T) {
		input := base64.StdEncoding.EncodeToString([]byte("padding test")) // Encodes "padding test" to base64
		expected := "padding test"
		result := DecodeBase64ToString(ctx, input)
		assert.Equal(t, types.StringValue(expected), result, "Should correctly decode a base64 string with padding")
	})
}
