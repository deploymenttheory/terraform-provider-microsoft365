package convert

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/stretchr/testify/assert"
)

// TestBitmaskEnum represents a mock bitmask enum for testing GraphToFrameworkBitmaskEnumAsSet
type TestBitmaskEnum int

const (
	TestBitmaskNone TestBitmaskEnum = 1
	TestBitmaskOne  TestBitmaskEnum = 2
	TestBitmaskTwo  TestBitmaskEnum = 4
	TestBitmaskAll  TestBitmaskEnum = 7
)

func (e TestBitmaskEnum) String() string {
	var values []string
	options := []string{"none", "one", "two"}
	for p := 0; p < 3; p++ {
		mask := TestBitmaskEnum(1 << p)
		if e&mask == mask {
			values = append(values, options[p])
		}
	}
	return strings.Join(values, ",")
}

func TestGraphToFrameworkString(t *testing.T) {
	// Case: Non-nil string pointer
	input := "test string"
	result := GraphToFrameworkString(&input)
	assert.Equal(t, types.StringValue(input), result, "Should return the string value")

	// Case: Nil string pointer
	var nilInput *string
	result = GraphToFrameworkString(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Empty string pointer
	emptyInput := ""
	result = GraphToFrameworkString(&emptyInput)
	assert.Equal(t, types.StringValue(""), result, "Should return empty string for empty string input")
}

func TestGraphToFrameworkTime(t *testing.T) {
	// Case: Nil time pointer
	var nilInput *time.Time
	result := GraphToFrameworkTime(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid time pointer
	input := time.Date(2023, 8, 8, 12, 0, 0, 0, time.UTC)
	expected := types.StringValue(input.Format(time.RFC3339))
	result = GraphToFrameworkTime(&input)
	assert.Equal(t, expected, result, "Should return the time formatted as RFC3339")

	// Case: Time with different location
	loc, _ := time.LoadLocation("America/New_York")
	locInput := time.Date(2023, 8, 8, 12, 0, 0, 0, loc)
	expected = types.StringValue(locInput.Format(time.RFC3339))
	result = GraphToFrameworkTime(&locInput)
	assert.Equal(t, expected, result, "Should return the time formatted as RFC3339 with correct timezone")
}

func TestGraphToFrameworkStringList(t *testing.T) {
	// Case: Nil input slice
	var nilInput []string
	result := GraphToFrameworkStringList(nilInput)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{}), result, "Should return empty list for nil input")

	// Case: Empty input slice
	emptyInput := []string{}
	result = GraphToFrameworkStringList(emptyInput)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{}), result, "Should return an empty list")

	// Case: Non-empty input slice
	input := []string{"one", "two", "three"}
	expected := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringValue("three"),
	}
	result = GraphToFrameworkStringList(input)
	assert.Equal(t, types.ListValueMust(types.StringType, expected), result, "Should convert slice of strings to list")
}

type mockEnum string

func (e mockEnum) String() string {
	return string(e)
}

func TestGraphToFrameworkEnumSlice(t *testing.T) {
	// Case: Nil input slice
	var nilInput []mockEnum
	result := GraphToFrameworkEnumSlice(nilInput)
	assert.Nil(t, result, "Should return nil for nil input")

	// Case: Empty input slice
	emptyInput := []mockEnum{}
	result = GraphToFrameworkEnumSlice(emptyInput)
	assert.Equal(t, 0, len(result), "Should return an empty slice")

	// Case: Non-empty input slice
	input := []mockEnum{"one", "two", "three"}
	expected := []types.String{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringValue("three"),
	}
	result = GraphToFrameworkEnumSlice(input)
	assert.Equal(t, expected, result, "Should convert slice of enums to slice of types.String")
}

func TestGraphToFrameworkBool(t *testing.T) {
	// Case: Nil bool pointer
	var nilInput *bool
	result := GraphToFrameworkBool(nilInput)
	assert.True(t, result.IsNull(), "Should return types.BoolNull() for nil input")

	// Case: True bool pointer
	trueInput := true
	result = GraphToFrameworkBool(&trueInput)
	assert.Equal(t, types.BoolValue(true), result, "Should return types.BoolValue(true) for true input")

	// Case: False bool pointer
	falseInput := false
	result = GraphToFrameworkBool(&falseInput)
	assert.Equal(t, types.BoolValue(false), result, "Should return types.BoolValue(false) for false input")
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

func TestGraphToFrameworkEnum(t *testing.T) {
	// Case: Nil enum pointer
	var nilInput *enumOneType
	result := GraphToFrameworkEnum(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid enum pointer
	input := EnumTwo
	result = GraphToFrameworkEnum(&input)
	assert.Equal(t, types.StringValue("Two"), result, "Should return the string representation of the enum")

	// Case: Different enum values
	testCases := []struct {
		input    enumOneType
		expected string
	}{
		{EnumOne, "One"},
		{EnumTwo, "Two"},
		{EnumThree, "Three"},
	}

	for _, tc := range testCases {
		result := GraphToFrameworkEnum(&tc.input)
		assert.Equal(t, types.StringValue(tc.expected), result, "Should return correct string for enum value")
	}
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

func TestGraphToFrameworkEnumPtrSlice(t *testing.T) {
	// Case: Nil input slice
	var nilInput []*enumTwoType
	result := GraphToFrameworkEnumPtrSlice(nilInput)
	assert.Nil(t, result, "Should return nil for nil input")

	// Case: Empty input slice
	emptyInput := []*enumTwoType{}
	result = GraphToFrameworkEnumPtrSlice(emptyInput)
	assert.Equal(t, 0, len(result), "Should return an empty slice")

	// Case: Non-empty input slice with valid enum pointers
	apple := EnumApple
	banana := EnumBanana
	input := []*enumTwoType{&apple, &banana}
	expected := []types.String{
		types.StringValue("Apple"),
		types.StringValue("Banana"),
	}
	result = GraphToFrameworkEnumPtrSlice(input)
	assert.Equal(t, expected, result, "Should convert slice of enum pointers to slice of types.String")

	// Case: Input slice with nil enum pointers
	apple = EnumApple
	input = []*enumTwoType{&apple, nil}
	expected = []types.String{
		types.StringValue("Apple"),
		types.StringNull(),
	}
	result = GraphToFrameworkEnumPtrSlice(input)
	assert.Equal(t, expected, result, "Should handle nil pointers correctly")
}

func TestGraphToFrameworkInt32AsInt64(t *testing.T) {
	// Case: Nil int32 pointer
	var nilInput *int32
	result := GraphToFrameworkInt32AsInt64(nilInput)
	assert.True(t, result.IsNull(), "Should return types.Int64Null() for nil input")

	// Case: Valid int32 pointer
	input := int32(42)
	result = GraphToFrameworkInt32AsInt64(&input)
	assert.Equal(t, types.Int64Value(42), result, "Should return types.Int64Value(42) for input 42")

	// Case: Negative int32 pointer
	negInput := int32(-123)
	result = GraphToFrameworkInt32AsInt64(&negInput)
	assert.Equal(t, types.Int64Value(-123), result, "Should return types.Int64Value(-123) for input -123")

	// Case: Max int32 pointer
	maxInput := int32(2147483647) // Max value for int32
	result = GraphToFrameworkInt32AsInt64(&maxInput)
	assert.Equal(t, types.Int64Value(2147483647), result, "Should correctly convert max int32 value")
}

func TestGraphToFrameworkInt32(t *testing.T) {
	// Case: Nil int32 pointer
	var nilInput *int32
	result := GraphToFrameworkInt32(nilInput)
	assert.True(t, result.IsNull(), "Should return types.Int32Null() for nil input")

	// Case: Valid int32 pointer
	input := int32(42)
	result = GraphToFrameworkInt32(&input)
	assert.Equal(t, types.Int32Value(42), result, "Should return types.Int32Value(42) for input 42")

	// Case: Negative int32 pointer
	negInput := int32(-123)
	result = GraphToFrameworkInt32(&negInput)
	assert.Equal(t, types.Int32Value(-123), result, "Should return types.Int32Value(-123) for input -123")

	// Case: Max int32 pointer
	maxInput := int32(2147483647) // Max value for int32
	result = GraphToFrameworkInt32(&maxInput)
	assert.Equal(t, types.Int32Value(2147483647), result, "Should correctly convert max int32 value")
}

func TestGraphToFrameworkDateOnly(t *testing.T) {
	// Case: Nil DateOnly pointer
	var nilInput *serialization.DateOnly
	result := GraphToFrameworkDateOnly(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid DateOnly pointer
	date := time.Date(2024, 8, 16, 0, 0, 0, 0, time.UTC)
	input := serialization.NewDateOnly(date)
	expected := types.StringValue("2024-08-16")
	result = GraphToFrameworkDateOnly(input)
	assert.Equal(t, expected, result, "Should return the date formatted as YYYY-MM-DD")

	// Case: Different valid DateOnly pointer
	date2 := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
	input2 := serialization.NewDateOnly(date2)
	expected2 := types.StringValue("1999-12-31")
	result = GraphToFrameworkDateOnly(input2)
	assert.Equal(t, expected2, result, "Should return the date formatted as YYYY-MM-DD")
}

func TestGraphToFrameworkBytes(t *testing.T) {
	// Case: Nil byte slice
	var nilInput []byte
	result := GraphToFrameworkBytes(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Empty byte slice
	emptyInput := []byte{}
	result = GraphToFrameworkBytes(emptyInput)
	assert.Equal(t, types.StringValue(""), result, "Should return empty string for empty byte slice")

	// Case: Non-empty byte slice
	input := []byte("Hello, World!")
	result = GraphToFrameworkBytes(input)
	assert.Equal(t, types.StringValue("Hello, World!"), result, "Should convert byte slice to string")
}

func TestGraphToFrameworkInt64(t *testing.T) {
	// Case: Nil int64 pointer
	var nilInput *int64
	result := GraphToFrameworkInt64(nilInput)
	assert.True(t, result.IsNull(), "Should return types.Int64Null() for nil input")

	// Case: Valid int64 pointer
	input := int64(42)
	result = GraphToFrameworkInt64(&input)
	assert.Equal(t, types.Int64Value(42), result, "Should return types.Int64Value(42) for input 42")

	// Case: Negative int64 pointer
	negInput := int64(-123)
	result = GraphToFrameworkInt64(&negInput)
	assert.Equal(t, types.Int64Value(-123), result, "Should return types.Int64Value(-123) for input -123")

	// Case: Max int64 pointer
	maxInput := int64(9223372036854775807) // Max value for int64
	result = GraphToFrameworkInt64(&maxInput)
	assert.Equal(t, types.Int64Value(9223372036854775807), result, "Should correctly convert max int64 value")
}

func TestGraphToFrameworkISODuration(t *testing.T) {
	// Case: Nil ISODuration pointer
	var nilInput *serialization.ISODuration
	result := GraphToFrameworkISODuration(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid ISODuration pointer
	input := serialization.NewDuration(1, 0, 3, 4, 5, 6, 7) // Example duration: P1Y3DT4H5M6S
	result = GraphToFrameworkISODuration(input)
	// Check that it contains the main components, as the exact string representation may vary
	assert.Contains(t, result.ValueString(), "P1Y3DT4H5M", "Should correctly convert valid ISODuration to ISO 8601 string")

	// Case: Valid ISODuration pointer with weeks
	weekInput := serialization.NewDuration(0, 2, 0, 0, 0, 0, 0) // Example duration: P2W
	result = GraphToFrameworkISODuration(weekInput)
	assert.Contains(t, result.ValueString(), "P2W", "Should correctly convert ISODuration with weeks to ISO 8601 string")
}

func TestGraphToFrameworkStringSet(t *testing.T) {
	ctx := context.Background()

	// Case: Empty slice
	result := GraphToFrameworkStringSet(ctx, []string{})
	assert.True(t, result.IsNull(), "Should return types.SetNull() for empty slice")

	// Case: Single string in slice
	input := []string{"one"}
	result = GraphToFrameworkStringSet(ctx, input)
	expected, _ := types.SetValueFrom(ctx, types.StringType, input)
	assert.Equal(t, expected, result, "Should return a Set with one element")

	// Case: Multiple strings in slice
	multiInput := []string{"a", "b", "c"}
	result = GraphToFrameworkStringSet(ctx, multiInput)
	expected, _ = types.SetValueFrom(ctx, types.StringType, multiInput)
	assert.Equal(t, expected, result, "Should return a Set with all input elements")
}

func TestGraphToFrameworkTimeOnly(t *testing.T) {
	// Case: Nil TimeOnly pointer
	var nilInput *serialization.TimeOnly
	result := GraphToFrameworkTimeOnly(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid TimeOnly pointer
	timeVal := time.Date(0, 1, 1, 14, 30, 45, 0, time.UTC)
	input := serialization.NewTimeOnly(timeVal)
	result = GraphToFrameworkTimeOnly(input)
	assert.Contains(t, result.ValueString(), "14:30:45", "Should return the time formatted correctly")
}

func TestGraphToFrameworkUUID(t *testing.T) {
	// Case: Nil UUID pointer
	var nilInput *uuid.UUID
	result := GraphToFrameworkUUID(nilInput)
	assert.True(t, result.IsNull(), "Should return types.StringNull() for nil input")

	// Case: Valid UUID pointer
	uuidVal := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	result = GraphToFrameworkUUID(&uuidVal)
	assert.Equal(t, types.StringValue("f47ac10b-58cc-4372-a567-0e02b2c3d479"), result, "Should return the UUID as a string")
}

func TestGraphToFrameworkStringSlice(t *testing.T) {
	// Case: Nil input slice
	var nilInput []string
	result := GraphToFrameworkStringSlice(nilInput)
	assert.Equal(t, 0, len(result), "Should return empty slice for nil input")

	// Case: Empty input slice
	emptyInput := []string{}
	result = GraphToFrameworkStringSlice(emptyInput)
	assert.Equal(t, 0, len(result), "Should return an empty slice")

	// Case: Non-empty input slice
	input := []string{"one", "two", "three"}
	expected := []types.String{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringValue("three"),
	}
	result = GraphToFrameworkStringSlice(input)
	assert.Equal(t, expected, result, "Should convert slice of strings to slice of types.String")
}

func TestGraphToFrameworkBitmaskEnumAsSet(t *testing.T) {
	ctx := context.Background()

	// Import MockBitmaskEnum from framework_to_graph_test.go

	// Case: Nil enum pointer
	var nilInput *MockBitmaskEnum
	result := GraphToFrameworkBitmaskEnumAsSet(ctx, nilInput)
	assert.True(t, result.IsNull(), "Should return null set for nil input")

	// Case: Empty string representation
	emptyEnum := MockBitmaskEnum(0)
	result = GraphToFrameworkBitmaskEnumAsSet(ctx, &emptyEnum)
	expected, _ := types.SetValueFrom(ctx, types.StringType, []string{})
	assert.Equal(t, expected, result, "Should return empty set for empty string representation")

	// Case: Single value enum
	singleEnum := MockBitmaskOne
	result = GraphToFrameworkBitmaskEnumAsSet(ctx, &singleEnum)
	expected, _ = types.SetValueFrom(ctx, types.StringType, []string{"one"})
	assert.Equal(t, expected, result, "Should return set with single value")

	// Case: Multiple values enum
	multiEnum := MockBitmaskOne | MockBitmaskTwo
	result = GraphToFrameworkBitmaskEnumAsSet(ctx, &multiEnum)
	// Order doesn't matter in a set, but we need to verify the values are correct
	elements := result.Elements()
	assert.Equal(t, 2, len(elements), "Should have 2 elements in the set")

	// Convert elements to strings for easier assertion
	var values []string
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			values = append(values, strVal.ValueString())
		}
	}

	assert.Contains(t, values, "one", "Set should contain 'one'")
	assert.Contains(t, values, "two", "Set should contain 'two'")

	// Case: All values enum
	allEnum := MockBitmaskAll
	result = GraphToFrameworkBitmaskEnumAsSet(ctx, &allEnum)
	elements = result.Elements()
	assert.Equal(t, 3, len(elements), "Should have 3 elements in the set")

	// Convert elements to strings for easier assertion
	values = []string{}
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			values = append(values, strVal.ValueString())
		}
	}

	assert.Contains(t, values, "none", "Set should contain 'none'")
	assert.Contains(t, values, "one", "Set should contain 'one'")
	assert.Contains(t, values, "two", "Set should contain 'two'")
}
