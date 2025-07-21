package convert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// ============================================================================
// GRAPH SDK → TERRAFORM FRAMEWORK (Read Operations)
// Used in Read() methods when populating Terraform state from Graph API responses
// ============================================================================

// GraphToFrameworkString converts a Graph SDK string pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
func GraphToFrameworkString(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(*value)
}

// GraphToFrameworkStringWithDefault converts a Graph SDK string pointer to a Terraform Framework string.
// Returns the default value if the pointer is nil or points to an empty string.
func GraphToFrameworkStringWithDefault(value *string, defaultValue string) types.String {
	if value == nil || *value == "" {
		return types.StringValue(defaultValue)
	}
	return types.StringValue(*value)
}

// GraphToFrameworkBool converts a Graph SDK bool pointer to a Terraform Framework bool.
// Returns types.BoolNull() if the input is nil.
func GraphToFrameworkBool(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

// GraphToFrameworkBoolWithDefault converts a Graph SDK bool pointer to a Terraform Framework bool.
// Returns the default value if the pointer is nil.
func GraphToFrameworkBoolWithDefault(value *bool, defaultValue bool) types.Bool {
	if value == nil {
		return types.BoolValue(defaultValue)
	}
	return types.BoolValue(*value)
}

// GraphToFrameworkInt32 converts a Graph SDK int32 pointer to a Terraform Framework int32.
// Returns types.Int32Null() if the input is nil.
func GraphToFrameworkInt32(value *int32) types.Int32 {
	if value == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*value)
}

// GraphToFrameworkInt64 converts a Graph SDK int64 pointer to a Terraform Framework int64.
// Returns types.Int64Null() if the input is nil.
func GraphToFrameworkInt64(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}

// GraphToFrameworkInt32AsInt64 converts a Graph SDK int32 pointer to a Terraform Framework int64.
// This is useful when the Graph API uses int32 but Terraform schema expects int64.
// Returns types.Int64Null() if the input is nil.
func GraphToFrameworkInt32AsInt64(value *int32) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

// GraphToFrameworkTime converts a Graph SDK time pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
// Time is formatted using RFC3339 format.
func GraphToFrameworkTime(value *time.Time) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.Format(constants.TimeFormatRFC3339Regex))
}

// GraphToFrameworkDateOnly converts a Graph SDK DateOnly pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
func GraphToFrameworkDateOnly(value *serialization.DateOnly) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.String())
}

// GraphToFrameworkTimeOnly converts a Graph SDK TimeOnly pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
// Uses default formatting with seconds precision (HH:MM:SS).
func GraphToFrameworkTimeOnly(value *serialization.TimeOnly) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.String())
}

// GraphToFrameworkTimeOnlyWithPrecision converts a Graph SDK TimeOnly pointer to a Terraform Framework string
// with the specified precision.
// Returns types.StringNull() if the input is nil.
// precision: 0-9, where 0 = HH:MM:SS, 1 = HH:MM:SS.f, 2 = HH:MM:SS.ff, etc.
func GraphToFrameworkTimeOnlyWithPrecision(value *serialization.TimeOnly, precision int) types.String {
	if value == nil {
		return types.StringNull()
	}

	// If precision is within valid range, use it
	if precision >= 0 && precision <= 9 {
		return types.StringValue(value.StringWithPrecision(precision))
	}

	// Otherwise use default precision
	return types.StringValue(value.String())
}

// GraphToFrameworkISODuration converts a Graph SDK ISODuration pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
// This function preserves the original ISO 8601 duration format as much as possible
// to avoid normalization issues (e.g., P7D becoming P1W) that can cause Terraform state inconsistencies.
func GraphToFrameworkISODuration(value *serialization.ISODuration) types.String {
	if value == nil {
		return types.StringNull()
	}

	// Reconstruct the ISO duration string manually to preserve the original format
	// This avoids the normalization that happens in ISODuration.String()
	var result string = "P"

	if value.GetYears() > 0 {
		result += fmt.Sprintf("%dY", value.GetYears())
	}

	if value.GetWeeks() > 0 {
		result += fmt.Sprintf("%dW", value.GetWeeks())
	}

	if value.GetDays() > 0 {
		result += fmt.Sprintf("%dD", value.GetDays())
	}

	// Add time component if needed
	if value.GetHours() > 0 || value.GetMinutes() > 0 || value.GetSeconds() > 0 || value.GetMilliSeconds() > 0 {
		result += "T"

		if value.GetHours() > 0 {
			result += fmt.Sprintf("%dH", value.GetHours())
		}

		if value.GetMinutes() > 0 {
			result += fmt.Sprintf("%dM", value.GetMinutes())
		}

		if value.GetSeconds() > 0 {
			result += fmt.Sprintf("%dS", value.GetSeconds())
		}

		// Milliseconds are typically not used in ISO 8601 durations in this context
		// but we'll handle them for completeness
		if value.GetMilliSeconds() > 0 {
			// If seconds are already present, append milliseconds as decimal
			if value.GetSeconds() > 0 {
				// Remove the S from the end
				result = result[:len(result)-1]
				result += fmt.Sprintf(".%03dS", value.GetMilliSeconds())
			} else {
				result += fmt.Sprintf("0.%03dS", value.GetMilliSeconds())
			}
		}
	}

	// Handle empty duration (just "P")
	if result == "P" {
		result = "PT0S"
	}

	return types.StringValue(result)
}

// GraphToFrameworkUUID converts a Graph SDK UUID pointer to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
func GraphToFrameworkUUID(value *uuid.UUID) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.String())
}

// GraphToFrameworkBytes converts a Graph SDK byte slice to a Terraform Framework string.
// Returns types.StringNull() if the input is nil.
// This is useful for script content which is stored as []byte but needs to be represented as a string.
func GraphToFrameworkBytes(value []byte) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(string(value))
}

// GraphToFrameworkEnum converts a Graph SDK enum pointer to a Terraform Framework string.
// Uses the String() method of the enum type to convert the value to a string.
// Returns types.StringNull() if the input is nil.
func GraphToFrameworkEnum[T fmt.Stringer](value *T) types.String {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue((*value).String())
}

// GraphToFrameworkStringList converts a Graph SDK string slice to a Terraform Framework list.
// Returns an empty list if the input is nil or empty.
func GraphToFrameworkStringList(value []string) types.List {
	if value == nil {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}

	values := make([]attr.Value, len(value))
	for i, v := range value {
		values[i] = types.StringValue(v)
	}

	return types.ListValueMust(types.StringType, values)
}

// GraphToFrameworkStringSet converts a Graph SDK string slice to a Terraform Framework set.
// Returns types.SetNull() if the input is empty.
func GraphToFrameworkStringSet(ctx context.Context, value []string) types.Set {
	if len(value) == 0 {
		return types.SetNull(types.StringType)
	}
	set, diags := types.SetValueFrom(ctx, types.StringType, value)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert string slice to types.Set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.StringType)
	}
	return set
}

// GraphToFrameworkStringSetPreserveEmpty converts a string slice to a Terraform Framework set,
// preserving empty arrays as empty sets instead of null.
// This is specifically needed for resources using HTTP clients (not Graph SDK) where we need
// to distinguish between:
// - Fields that were configured in Terraform as empty arrays [] → should become empty sets
// - Fields that were not configured in Terraform at all → should remain null
//
// Use this helper when you need to maintain state consistency for optional set attributes
// that can be explicitly set to empty in Terraform configuration.
func GraphToFrameworkStringSetPreserveEmpty(ctx context.Context, value []string) types.Set {
	set, diags := types.SetValueFrom(ctx, types.StringType, value)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert string slice to types.Set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.StringType)
	}
	return set
}

// GraphToFrameworkEnumSlice converts a Graph SDK enum slice to a slice of Terraform Framework strings.
// Uses the String() method of the enum type to convert each value to a string.
// Returns nil if the input is nil.
func GraphToFrameworkEnumSlice[T fmt.Stringer](value []T) []types.String {
	if value == nil {
		return nil
	}

	result := make([]types.String, len(value))
	for i, v := range value {
		result[i] = types.StringValue(v.String())
	}

	return result
}

// GraphToFrameworkEnumPtrSlice converts a Graph SDK enum pointer slice to a slice of Terraform Framework strings.
// Uses the String() method of the enum type to convert each value to a string.
// Returns types.StringNull() for nil pointers in the slice.
func GraphToFrameworkEnumPtrSlice[T fmt.Stringer](value []*T) []types.String {
	if value == nil {
		return nil
	}

	result := make([]types.String, len(value))
	for i, v := range value {
		if v == nil {
			result[i] = types.StringNull()
		} else {
			result[i] = types.StringValue((*v).String())
		}
	}

	return result
}

// GraphToFrameworkStringSlice converts a Graph SDK string slice to a slice of Terraform Framework strings.
// Returns an empty slice if the input is nil or empty.
func GraphToFrameworkStringSlice(value []string) []types.String {
	if value == nil {
		return []types.String{}
	}

	result := make([]types.String, len(value))
	for i, v := range value {
		result[i] = types.StringValue(v)
	}

	return result
}

// GraphToFrameworkFloat64 converts a Graph SDK float64 pointer to a Terraform Framework float64.
// Returns types.Float64Null() if the input is nil.
func GraphToFrameworkFloat64(value *float64) types.Float64 {
	if value == nil {
		return types.Float64Null()
	}
	return types.Float64Value(*value)
}

// GraphToFrameworkBitmaskEnumAsSet converts a bitmask enum to a Terraform Framework set of strings.
// This is useful for APIs that use bitmask enums with String() methods that return comma-separated values.
// The function calls the String() method on the enum, splits the result by commas, and creates a set.
// Returns types.SetNull() if the input is nil.
func GraphToFrameworkBitmaskEnumAsSet[T fmt.Stringer](ctx context.Context, value *T) types.Set {
	if value == nil {
		return types.SetNull(types.StringType)
	}

	// Get string representation and split by commas
	enumStr := (*value).String()
	if enumStr == "" {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Split the comma-separated string
	parts := strings.Split(enumStr, ",")

	// Create a set from the parts
	set, diags := types.SetValueFrom(ctx, types.StringType, parts)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert bitmask enum to set", map[string]interface{}{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.StringType)
	}

	return set
}
