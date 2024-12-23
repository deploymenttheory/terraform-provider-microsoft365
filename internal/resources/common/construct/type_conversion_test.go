package construct

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSetStringProperty(t *testing.T) {
	var result *string

	// Case: Valid string
	result = nil
	optString := types.StringValue("test")
	SetStringProperty(optString, func(val *string) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, "test", *result)

	// Case: Null value
	result = nil
	optString = types.StringNull()
	SetStringProperty(optString, func(val *string) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown()
	SetStringProperty(optString, func(val *string) {
		result = val
	})
	assert.Nil(t, result)
}

func TestSetBoolProperty(t *testing.T) {
	var result *bool

	// Case: Valid bool
	result = nil
	optBool := types.BoolValue(true)
	SetBoolProperty(optBool, func(val *bool) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, true, *result)

	// Case: Null value
	result = nil
	optBool = types.BoolNull()
	SetBoolProperty(optBool, func(val *bool) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optBool = types.BoolUnknown()
	SetBoolProperty(optBool, func(val *bool) {
		result = val
	})
	assert.Nil(t, result)
}

func TestSetInt32Property(t *testing.T) {
	var result *int32

	// Case: Valid int32
	result = nil
	optInt := types.Int32Value(123)
	SetInt32Property(optInt, func(val *int32) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, int32(123), *result)

	// Case: Null value
	result = nil
	optInt = types.Int32Null()
	SetInt32Property(optInt, func(val *int32) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optInt = types.Int32Unknown()
	SetInt32Property(optInt, func(val *int32) {
		result = val
	})
	assert.Nil(t, result)
}

func TestSetInt64Property(t *testing.T) {
	var result *int64

	// Case: Valid int64
	result = nil
	optInt := types.Int64Value(456)
	SetInt64Property(optInt, func(val *int64) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, int64(456), *result)

	// Case: Null value
	result = nil
	optInt = types.Int64Null()
	SetInt64Property(optInt, func(val *int64) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optInt = types.Int64Unknown()
	SetInt64Property(optInt, func(val *int64) {
		result = val
	})
	assert.Nil(t, result)
}

func TestParseEnum(t *testing.T) {
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
	err := ParseEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "parsed", *result)

	// Case: Invalid enum value
	result = nil
	optEnum = types.StringValue("invalid") // Simulate an invalid enum value
	err = ParseEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.Error(t, err)
	assert.Nil(t, result)

	// Case: Null value
	result = nil
	optEnum = types.StringNull() // Simulate a null StringValue
	err = ParseEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optEnum = types.StringUnknown() // Simulate an unknown StringValue
	err = ParseEnum[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSetArrayProperty(t *testing.T) {
	var result []string

	// Case: Valid array with all valid strings
	result = nil
	validArray := []types.String{
		types.StringValue("test1"),
		types.StringValue("test2"),
		types.StringValue("test3"),
	}
	SetArrayProperty(validArray, func(val []string) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, []string{"test1", "test2", "test3"}, result)

	// Case: Array with mix of valid, null, and unknown values
	result = nil
	mixedArray := []types.String{
		types.StringValue("test1"),
		types.StringNull(),
		types.StringValue("test3"),
		types.StringUnknown(),
	}
	SetArrayProperty(mixedArray, func(val []string) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, []string{"test1", "test3"}, result)

	// Case: Empty array
	result = nil
	emptyArray := []types.String{}
	SetArrayProperty(emptyArray, func(val []string) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Array with only null and unknown values
	result = nil
	nullUnknownArray := []types.String{
		types.StringNull(),
		types.StringUnknown(),
		types.StringNull(),
	}
	SetArrayProperty(nullUnknownArray, func(val []string) {
		result = val
	})
	assert.Nil(t, result)

	// Case: Array with one valid value
	result = nil
	singleValidArray := []types.String{
		types.StringValue("test"),
	}
	SetArrayProperty(singleValidArray, func(val []string) {
		result = val
	})
	assert.NotNil(t, result)
	assert.Equal(t, []string{"test"}, result)

	// Case: Nil array
	result = nil
	var nilArray []types.String
	SetArrayProperty(nilArray, func(val []string) {
		result = val
	})
	assert.Nil(t, result)
}
