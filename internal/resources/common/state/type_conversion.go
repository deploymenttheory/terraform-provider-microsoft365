package state

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// StringPtrToString converts a string pointer to a string.
func StringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// TimeToString converts a *time.Time to a types.String.
// If the input is nil, it returns types.StringNull().
// Otherwise, it returns a types.String with the time formatted in RFC3339 format.
func TimeToString(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(helpers.TimeFormatRFC3339))
}

// SliceToTypeStringSlice converts a slice of strings to a slice of types.String.
// It handles nil input by returning nil, and empty slices by returning an empty slice of types.String.
func SliceToTypeStringSlice(input []string) []types.String {
	if input == nil {
		return nil
	}

	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}

	return result
}

// EnumSliceToTypeStringSlice converts a slice of enum-like constants to a slice of types.String.
// It uses the String() method of the enum type to convert each value to a string.
func EnumSliceToTypeStringSlice[T fmt.Stringer](input []T) []types.String {
	if input == nil {
		return nil
	}

	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v.String())
	}

	return result
}

// BoolPtrToTypeBool converts a *bool to a types.Bool.
// If the input is nil, it returns types.BoolNull().
func BoolPtrToTypeBool(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// EnumPtrToTypeString converts a pointer to an enum-like type to a types.String.
// It uses the String() method of the enum type to convert the value to a string.
func EnumPtrToTypeString[T fmt.Stringer](e *T) types.String {
	if e == nil {
		return types.StringNull()
	}
	return types.StringValue((*e).String())
}

// EnumListPtrToTypeStringSlice converts a slice of pointers to enum-like constants to a slice of types.String.
// It uses the String() method of the enum type to convert each value to a string.
func EnumListPtrToTypeStringSlice[T fmt.Stringer](input []*T) []types.String {
	if input == nil {
		return nil
	}

	result := make([]types.String, len(input))
	for i, v := range input {
		if v == nil {
			result[i] = types.StringNull()
		} else {
			// Dereference the pointer and call String()
			result[i] = types.StringValue((*v).String())
		}
	}

	return result
}

// Int32PtrToTypeInt64 converts a pointer to an int32 to a types.Int64.
// This function is useful for converting nullable int32 values from the SDK to Terraform's types.Int64.
func Int32PtrToTypeInt64(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*i))
}

// Int32PtrToTypeInt32 converts a pointer to an int32 to a types.Int32.
// This function is useful for converting nullable int32 values from the SDK to Terraform's types.Int32.
func Int32PtrToTypeInt32(i *int32) types.Int32 {
	if i == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*i)
}

// DateOnlyPtrToString converts a DateOnly pointer to a Terraform string.
func DateOnlyPtrToString(date *serialization.DateOnly) types.String {
	if date == nil {
		return types.StringNull()
	}
	return types.StringValue(date.String())
}

// ByteToString converts a byte slice to a string.
// It returns the byte slice encoded as a base64 string.
func ByteToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// BoolPtrToBool converts a bool pointer to a bool.
// If the input is nil, it returns false.
func BoolPtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Int64PtrToTypeInt64 converts a *int64 to a types.Int64.
// If the input is nil, it returns types.Int64Null().
func Int64PtrToTypeInt64(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*i)
}

// ISO8601DurationToString converts an ISODuration to a types.String value.
func ISO8601DurationToString(duration *serialization.ISODuration) types.String {
	if duration == nil {
		return types.StringNull()
	}
	return types.StringValue(duration.String())
}

// DecodeBase64ToString decodes a base64-encoded string and returns a basetypes.StringValue.
// If decoding fails, it logs a warning and returns the original string as a basetypes.StringValue.
func DecodeBase64ToString(ctx context.Context, encoded string) types.String {
	decodedContent, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		tflog.Warn(ctx, "Failed to decode base64 content", map[string]interface{}{
			"error": err.Error(),
		})
		// Return the original string as a fallback
		return types.StringValue(encoded)
	}

	// Return the decoded content as types.StringValue
	return types.StringValue(string(decodedContent))
}

// StringListToTypeList converts a slice of strings to a types.List.
func StringListToTypeList(strings []string) types.List {
	values := make([]attr.Value, len(strings))
	for i, s := range strings {
		values[i] = types.StringValue(s)
	}

	list, _ := basetypes.NewListValue(types.StringType, values)
	return list
}
