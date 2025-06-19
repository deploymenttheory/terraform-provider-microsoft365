package convert

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/stretchr/testify/assert"
)

func TestFrameworkToGraphString(t *testing.T) {
	var result *string

	// Case: Valid string
	result = nil
	optString := types.StringValue("test")
	FrameworkToGraphString(optString, func(val *string) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, "test", *result)

	// Case: Null value
	result = nil
	optString = types.StringNull()
	FrameworkToGraphString(optString, func(val *string) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown()
	FrameworkToGraphString(optString, func(val *string) {
		result = val
	})
	assert.Nil(t, result)
}

func TestFrameworkToGraphBool(t *testing.T) {
	var result *bool

	// Case: Valid bool
	result = nil
	optBool := types.BoolValue(true)
	FrameworkToGraphBool(optBool, func(val *bool) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, true, *result)

	// Case: Null value
	result = nil
	optBool = types.BoolNull()
	FrameworkToGraphBool(optBool, func(val *bool) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optBool = types.BoolUnknown()
	FrameworkToGraphBool(optBool, func(val *bool) {
		result = val
	})
	assert.Nil(t, result)
}

func TestFrameworkToGraphInt32(t *testing.T) {
	var result *int32

	// Case: Valid int32
	result = nil
	optInt := types.Int32Value(123)
	FrameworkToGraphInt32(optInt, func(val *int32) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, int32(123), *result)

	// Case: Null value
	result = nil
	optInt = types.Int32Null()
	FrameworkToGraphInt32(optInt, func(val *int32) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optInt = types.Int32Unknown()
	FrameworkToGraphInt32(optInt, func(val *int32) {
		result = val
	})
	assert.Nil(t, result)
}

func TestFrameworkToGraphInt64(t *testing.T) {
	var result *int64

	// Case: Valid int64
	result = nil
	optInt := types.Int64Value(456)
	FrameworkToGraphInt64(optInt, func(val *int64) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, int64(456), *result)

	// Case: Null value
	result = nil
	optInt = types.Int64Null()
	FrameworkToGraphInt64(optInt, func(val *int64) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optInt = types.Int64Unknown()
	FrameworkToGraphInt64(optInt, func(val *int64) {
		result = val
	})
	assert.Nil(t, result)
}

func TestFrameworkToGraphEnum(t *testing.T) {
	var result *string

	// Parser that returns the expected type `any`, simulating a valid string enum parsing result.
	parser := func(input string) (any, error) {
		if input == "valid" {
			return "parsed", nil // Simulate a valid enum parsing result as `any`
		}
		return nil, errors.New("invalid value") // Simulate a parsing error for invalid input
	}

	// Case: Valid enum value
	result = nil
	optEnum := types.StringValue("valid") // Simulate a valid StringValue from Terraform SDK
	err := FrameworkToGraphEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "parsed", *result)

	// Case: Invalid enum value
	result = nil
	optEnum = types.StringValue("invalid") // Simulate an invalid enum value
	err = FrameworkToGraphEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.Error(t, err)
	assert.Nil(t, result)

	// Case: Null value
	result = nil
	optEnum = types.StringNull() // Simulate a null StringValue
	err = FrameworkToGraphEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optEnum = types.StringUnknown() // Simulate an unknown StringValue
	err = FrameworkToGraphEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFrameworkToGraphBytes(t *testing.T) {
	var result []byte

	// Case: Valid string value
	result = nil
	optString := types.StringValue("test content")
	FrameworkToGraphBytes(optString, func(val []byte) {
		result = val
	})
	assert.NotNil(t, result, "Setter should be called for a valid string value")
	assert.Equal(t, []byte("test content"), result, "Setter should receive the correct byte slice")

	// Case: Empty string value
	result = nil
	optString = types.StringValue("")
	FrameworkToGraphBytes(optString, func(val []byte) {
		result = val
	})
	assert.NotNil(t, result, "Setter should be called for an empty string value")
	assert.Equal(t, []byte(""), result, "Setter should receive an empty byte slice")

	// Case: Null value
	result = nil
	optString = types.StringNull()
	FrameworkToGraphBytes(optString, func(val []byte) {
		result = val
	})
	assert.Nil(t, result, "Setter should not be called for a null value")

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown()
	FrameworkToGraphBytes(optString, func(val []byte) {
		result = val
	})
	assert.Nil(t, result, "Setter should not be called for an unknown value")
}

func TestFrameworkToGraphISODuration(t *testing.T) {
	var result *serialization.ISODuration

	// Case: Valid ISO 8601 duration
	result = nil
	optString := types.StringValue("P1Y2M3DT4H5M6S") // Valid ISO 8601 duration
	err := FrameworkToGraphISODuration(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for valid ISO 8601 duration")
	assert.NotNil(t, result, "Setter should be called for a valid ISO 8601 duration")
	// The ISODuration.String() method may normalize the duration and omit seconds
	// So we check that the string contains the main components
	assert.Contains(t, result.String(), "P1Y2M3DT4H5M")

	// Case: Invalid ISO 8601 duration
	result = nil
	optString = types.StringValue("InvalidDuration") // Invalid ISO 8601 duration
	err = FrameworkToGraphISODuration(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.Error(t, err, "An error should occur for invalid ISO 8601 duration")
	assert.Nil(t, result, "Setter should not be called for an invalid ISO 8601 duration")

	// Case: Null value
	result = nil
	optString = types.StringNull() // Null value
	err = FrameworkToGraphISODuration(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for a null value")
	assert.Nil(t, result, "Setter should not be called for a null value")

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown() // Unknown value
	err = FrameworkToGraphISODuration(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for an unknown value")
	assert.Nil(t, result, "Setter should not be called for an unknown value")
}

func TestFrameworkToGraphObjectsFromStringSet(t *testing.T) {
	// Define a test type and converter function
	type TestObject struct {
		ID   string
		Name string
	}

	converter := func(ctx context.Context, values []string) []TestObject {
		result := make([]TestObject, 0, len(values))
		for _, val := range values {
			result = append(result, TestObject{
				ID:   val + "_id",
				Name: val,
			})
		}
		return result
	}

	// Case: Valid string set
	var result []TestObject
	elements := []attr.Value{
		types.StringValue("test1"),
		types.StringValue("test2"),
		types.StringValue("test3"),
	}
	set, diags := types.SetValue(types.StringType, elements)
	assert.False(t, diags.HasError())

	err := FrameworkToGraphObjectsFromStringSet(context.Background(), set, converter, func(val []TestObject) {
		result = val
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "test1_id", result[0].ID)
	assert.Equal(t, "test1", result[0].Name)
	assert.Equal(t, "test2_id", result[1].ID)
	assert.Equal(t, "test2", result[1].Name)
	assert.Equal(t, "test3_id", result[2].ID)
	assert.Equal(t, "test3", result[2].Name)

	// Case: Empty string set
	result = nil
	emptySet, diags := types.SetValue(types.StringType, []attr.Value{})
	assert.False(t, diags.HasError())

	err = FrameworkToGraphObjectsFromStringSet(context.Background(), emptySet, converter, func(val []TestObject) {
		result = val
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))

	// Case: Null set
	result = nil
	nullSet := types.SetNull(types.StringType)

	err = FrameworkToGraphObjectsFromStringSet(context.Background(), nullSet, converter, func(val []TestObject) {
		result = val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown set
	result = nil
	unknownSet := types.SetUnknown(types.StringType)

	err = FrameworkToGraphObjectsFromStringSet(context.Background(), unknownSet, converter, func(val []TestObject) {
		result = val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFrameworkToGraphTime(t *testing.T) {
	// Test case: Valid RFC3339 string
	validTimeStr := "2023-01-15T08:30:00Z"
	var resultTime *time.Time

	err := FrameworkToGraphTime(types.StringValue(validTimeStr), func(t *time.Time) {
		resultTime = t
	})

	assert.NoError(t, err)
	assert.NotNil(t, resultTime)
	expectedTime, _ := time.Parse(time.RFC3339, validTimeStr)
	assert.Equal(t, expectedTime.UTC(), resultTime.UTC())

	// Test case: Invalid time string
	invalidTimeStr := "not-a-valid-time"
	err = FrameworkToGraphTime(types.StringValue(invalidTimeStr), func(t *time.Time) {
		resultTime = t
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse time string")

	// Test case: Empty string
	resultTime = nil
	err = FrameworkToGraphTime(types.StringValue(""), func(t *time.Time) {
		resultTime = t
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTime)

	// Test case: Null string
	resultTime = nil
	err = FrameworkToGraphTime(types.StringNull(), func(t *time.Time) {
		resultTime = t
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTime)

	// Test case: Unknown string
	resultTime = nil
	err = FrameworkToGraphTime(types.StringUnknown(), func(t *time.Time) {
		resultTime = t
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTime)
}

func TestFrameworkToGraphTimeOnly(t *testing.T) {
	// Test case: Valid time string
	validTimeStr := "08:30:00"
	var resultTimeOnly *serialization.TimeOnly

	err := FrameworkToGraphTimeOnly(types.StringValue(validTimeStr), func(to *serialization.TimeOnly) {
		resultTimeOnly = to
	})

	assert.NoError(t, err)
	assert.NotNil(t, resultTimeOnly)
	assert.Contains(t, resultTimeOnly.String(), validTimeStr)

	// Test case: Invalid time string
	invalidTimeStr := "not-a-valid-time"
	err = FrameworkToGraphTimeOnly(types.StringValue(invalidTimeStr), func(to *serialization.TimeOnly) {
		resultTimeOnly = to
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse time string")

	// Test case: Empty string
	resultTimeOnly = nil
	err = FrameworkToGraphTimeOnly(types.StringValue(""), func(to *serialization.TimeOnly) {
		resultTimeOnly = to
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTimeOnly)

	// Test case: Null string
	resultTimeOnly = nil
	err = FrameworkToGraphTimeOnly(types.StringNull(), func(to *serialization.TimeOnly) {
		resultTimeOnly = to
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTimeOnly)

	// Test case: Unknown string
	resultTimeOnly = nil
	err = FrameworkToGraphTimeOnly(types.StringUnknown(), func(to *serialization.TimeOnly) {
		resultTimeOnly = to
	})

	assert.NoError(t, err)
	assert.Nil(t, resultTimeOnly)

	// Test case: Boundary values
	boundaryTimes := []string{
		"00:00:00", // Midnight
		"23:59:59", // Just before midnight
		"12:00:00", // Noon
	}

	for _, timeStr := range boundaryTimes {
		err = FrameworkToGraphTimeOnly(types.StringValue(timeStr), func(to *serialization.TimeOnly) {
			resultTimeOnly = to
		})

		assert.NoError(t, err)
		assert.NotNil(t, resultTimeOnly)
		assert.Contains(t, resultTimeOnly.String(), timeStr[:8]) // First 8 chars (HH:MM:SS)
	}
}

func TestFrameworkToGraphUUID(t *testing.T) {
	expectedUUID := uuid.New()
	val := basetypes.NewStringValue(expectedUUID.String())

	var actual *uuid.UUID
	setter := func(u *uuid.UUID) {
		actual = u
	}

	err := FrameworkToGraphUUID(val, setter)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expectedUUID, *actual)

	// Case: Null value
	val = basetypes.NewStringNull()
	actual = nil
	err = FrameworkToGraphUUID(val, setter)

	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Case: Unknown value
	val = basetypes.NewStringUnknown()
	actual = nil
	err = FrameworkToGraphUUID(val, setter)

	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Case: Empty string
	val = basetypes.NewStringValue("")
	actual = nil
	err = FrameworkToGraphUUID(val, setter)

	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Case: Invalid UUID
	val = basetypes.NewStringValue("not-a-uuid")
	actual = nil
	err = FrameworkToGraphUUID(val, setter)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestFrameworkToGraphStringList(t *testing.T) {
	ctx := context.Background()

	// Case: Valid list
	var result []string
	elements := []attr.Value{
		types.StringValue("value1"),
		types.StringValue("value2"),
		types.StringValue("value3"),
	}
	list, diags := types.ListValue(types.StringType, elements)
	assert.False(t, diags.HasError())

	err := FrameworkToGraphStringList(ctx, list, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{"value1", "value2", "value3"}, result)

	// Case: Empty list
	result = nil
	emptyList, diags := types.ListValue(types.StringType, []attr.Value{})
	assert.False(t, diags.HasError())

	err = FrameworkToGraphStringList(ctx, emptyList, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{}, result)

	// Case: Null list
	result = nil
	nullList := types.ListNull(types.StringType)

	err = FrameworkToGraphStringList(ctx, nullList, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown list
	result = nil
	unknownList := types.ListUnknown(types.StringType)

	err = FrameworkToGraphStringList(ctx, unknownList, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFrameworkToGraphStringSet(t *testing.T) {
	ctx := context.Background()

	// Case: Valid set
	var result []string
	elements := []attr.Value{
		types.StringValue("value1"),
		types.StringValue("value2"),
		types.StringValue("value3"),
	}
	set, diags := types.SetValue(types.StringType, elements)
	assert.False(t, diags.HasError())

	err := FrameworkToGraphStringSet(ctx, set, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"value1", "value2", "value3"}, result)

	// Case: Empty set
	result = nil
	emptySet, diags := types.SetValue(types.StringType, []attr.Value{})
	assert.False(t, diags.HasError())

	err = FrameworkToGraphStringSet(ctx, emptySet, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{}, result)

	// Case: Null set
	result = nil
	nullSet := types.SetNull(types.StringType)

	err = FrameworkToGraphStringSet(ctx, nullSet, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown set
	result = nil
	unknownSet := types.SetUnknown(types.StringType)

	err = FrameworkToGraphStringSet(ctx, unknownSet, func(val []string) {
		result = val
	})

	assert.NoError(t, err)
	assert.Nil(t, result)
}

// --- MOCK ENUM ---

type MockBrandingOptions int

const (
	MOCK_NONE    MockBrandingOptions = 1
	MOCK_LOGO    MockBrandingOptions = 2
	MOCK_NAME    MockBrandingOptions = 4
	MOCK_CONTACT MockBrandingOptions = 8
)

func MockParseBrandingOptions(input string) (any, error) {
	var result MockBrandingOptions
	parts := strings.Split(input, ",")
	for _, str := range parts {
		switch strings.TrimSpace(str) {
		case "none":
			result |= MOCK_NONE
		case "logo":
			result |= MOCK_LOGO
		case "name":
			result |= MOCK_NAME
		case "contact":
			result |= MOCK_CONTACT
		default:
			return nil, nil // simulate Microsoft behavior
		}
	}
	return &result, nil
}

// --- TESTS ---

func TestFrameworkToGraphBitmaskEnum(t *testing.T) {
	// Case: Valid single value
	val := basetypes.NewStringValue("logo")
	var actual *MockBrandingOptions
	setter := func(e *MockBrandingOptions) {
		actual = e
	}

	err := FrameworkToGraphBitmaskEnum(val, MockParseBrandingOptions, setter)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, MOCK_LOGO, *actual)

	// Case: Valid multiple values
	val = basetypes.NewStringValue("logo,name")
	actual = nil
	err = FrameworkToGraphBitmaskEnum(val, MockParseBrandingOptions, setter)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	expected := MOCK_LOGO | MOCK_NAME
	assert.Equal(t, expected, *actual)

	// Case: Invalid value
	val = basetypes.NewStringValue("invalid")
	actual = nil
	err = FrameworkToGraphBitmaskEnum(val, MockParseBrandingOptions, setter)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Case: Null value
	val = basetypes.NewStringNull()
	actual = nil
	err = FrameworkToGraphBitmaskEnum(val, MockParseBrandingOptions, setter)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Case: Unknown value
	val = basetypes.NewStringUnknown()
	actual = nil
	err = FrameworkToGraphBitmaskEnum(val, MockParseBrandingOptions, setter)
	assert.NoError(t, err)
	assert.Nil(t, actual)
}

func TestFrameworkToGraphDateOnly(t *testing.T) {
	// Test case: Valid date string
	validDateStr := "2023-01-15"
	var resultDate *serialization.DateOnly

	err := FrameworkToGraphDateOnly(types.StringValue(validDateStr), func(d *serialization.DateOnly) {
		resultDate = d
	})

	assert.NoError(t, err)
	assert.NotNil(t, resultDate)
	assert.Equal(t, validDateStr, resultDate.String())

	// Test case: Invalid date string
	invalidDateStr := "not-a-valid-date"
	err = FrameworkToGraphDateOnly(types.StringValue(invalidDateStr), func(d *serialization.DateOnly) {
		resultDate = d
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse date string")

	// Test case: Empty string
	resultDate = nil
	err = FrameworkToGraphDateOnly(types.StringValue(""), func(d *serialization.DateOnly) {
		resultDate = d
	})

	assert.NoError(t, err)
	assert.Nil(t, resultDate)

	// Test case: Null string
	resultDate = nil
	err = FrameworkToGraphDateOnly(types.StringNull(), func(d *serialization.DateOnly) {
		resultDate = d
	})

	assert.NoError(t, err)
	assert.Nil(t, resultDate)

	// Test case: Unknown string
	resultDate = nil
	err = FrameworkToGraphDateOnly(types.StringUnknown(), func(d *serialization.DateOnly) {
		resultDate = d
	})

	assert.NoError(t, err)
	assert.Nil(t, resultDate)
}
