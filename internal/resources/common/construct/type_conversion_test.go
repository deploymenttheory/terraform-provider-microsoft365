package construct

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
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

func TestSetEnumProperty(t *testing.T) {
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
	err := SetEnumProperty[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "parsed", *result)

	// Case: Invalid enum value
	result = nil
	optEnum = types.StringValue("invalid") // Simulate an invalid enum value
	err = SetEnumProperty[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.Error(t, err)
	assert.Nil(t, result)

	// Case: Null value
	result = nil
	optEnum = types.StringNull() // Simulate a null StringValue
	err = SetEnumProperty[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Case: Unknown value
	result = nil
	optEnum = types.StringUnknown() // Simulate an unknown StringValue
	err = SetEnumProperty[string](optEnum, parser, func(val string) {
		result = &val
	})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSetBytesProperty(t *testing.T) {
	var result []byte

	// Case: Valid string value
	result = nil
	optString := types.StringValue("test content")
	SetBytesProperty(optString, func(val []byte) {
		result = val
	})
	assert.NotNil(t, result, "Setter should be called for a valid string value")
	assert.Equal(t, []byte("test content"), result, "Setter should receive the correct byte slice")

	// Case: Empty string value
	result = nil
	optString = types.StringValue("")
	SetBytesProperty(optString, func(val []byte) {
		result = val
	})
	assert.NotNil(t, result, "Setter should be called for an empty string value")
	assert.Equal(t, []byte(""), result, "Setter should receive an empty byte slice")

	// Case: Null value
	result = nil
	optString = types.StringNull()
	SetBytesProperty(optString, func(val []byte) {
		result = val
	})
	assert.Nil(t, result, "Setter should not be called for a null value")

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown()
	SetBytesProperty(optString, func(val []byte) {
		result = val
	})
	assert.Nil(t, result, "Setter should not be called for an unknown value")
}

func TestSetISODurationProperty(t *testing.T) {
	var result *serialization.ISODuration

	// Case: Valid ISO 8601 duration
	result = nil
	optString := types.StringValue("P1Y2M3DT4H5M6S") // Valid ISO 8601 duration
	err := SetISODurationProperty(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for valid ISO 8601 duration")
	assert.NotNil(t, result, "Setter should be called for a valid ISO 8601 duration")
	assert.Equal(t, "P1Y2M3DT4H5M", result.String(), "The parsed duration should match the normalized representation")

	// Case: Invalid ISO 8601 duration
	result = nil
	optString = types.StringValue("InvalidDuration") // Invalid ISO 8601 duration
	err = SetISODurationProperty(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.Error(t, err, "An error should occur for invalid ISO 8601 duration")
	assert.Nil(t, result, "Setter should not be called for an invalid ISO 8601 duration")

	// Case: Null value
	result = nil
	optString = types.StringNull() // Null value
	err = SetISODurationProperty(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for a null value")
	assert.Nil(t, result, "Setter should not be called for a null value")

	// Case: Unknown value
	result = nil
	optString = types.StringUnknown() // Unknown value
	err = SetISODurationProperty(optString, func(val *serialization.ISODuration) {
		result = val
	})
	assert.NoError(t, err, "No error should occur for an unknown value")
	assert.Nil(t, result, "Setter should not be called for an unknown value")
}
