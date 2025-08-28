package convert

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ============================================================================
// MAP → TERRAFORM FRAMEWORK (Read Operations)
// Used in Read() methods when receiving API responses as map[string]interface{}
// ============================================================================

// MapToFrameworkString extracts a string value from a map and converts it to a Terraform Framework string.
// Returns null if the key doesn't exist or the value is not a string.
func MapToFrameworkString(data map[string]any, key string) types.String {
	if value, ok := data[key].(string); ok {
		return types.StringValue(value)
	}
	return types.StringNull()
}

// MapToFrameworkBool extracts a bool value from a map and converts it to a Terraform Framework bool.
// Returns null if the key doesn't exist or the value is not a bool.
func MapToFrameworkBool(data map[string]any, key string) types.Bool {
	if value, ok := data[key].(bool); ok {
		return types.BoolValue(value)
	}
	return types.BoolNull()
}

// MapToFrameworkInt32 extracts an int32 value from a map and converts it to a Terraform Framework int32.
// Returns null if the key doesn't exist or the value is not convertible to int32.
func MapToFrameworkInt32(data map[string]any, key string) types.Int32 {
	switch value := data[key].(type) {
	case int32:
		return types.Int32Value(value)
	case int:
		return types.Int32Value(int32(value))
	case float64:
		return types.Int32Value(int32(value))
	}
	return types.Int32Null()
}

// MapToFrameworkInt64 extracts an int64 value from a map and converts it to a Terraform Framework int64.
// Returns null if the key doesn't exist or the value is not convertible to int64.
func MapToFrameworkInt64(data map[string]any, key string) types.Int64 {
	switch value := data[key].(type) {
	case int64:
		return types.Int64Value(value)
	case int:
		return types.Int64Value(int64(value))
	case int32:
		return types.Int64Value(int64(value))
	case float64:
		return types.Int64Value(int64(value))
	}
	return types.Int64Null()
}

// MapToFrameworkFloat64 extracts a float64 value from a map and converts it to a Terraform Framework float64.
// Returns null if the key doesn't exist or the value is not convertible to float64.
func MapToFrameworkFloat64(data map[string]any, key string) types.Float64 {
	switch value := data[key].(type) {
	case float64:
		return types.Float64Value(value)
	case float32:
		return types.Float64Value(float64(value))
	case int:
		return types.Float64Value(float64(value))
	case int32:
		return types.Float64Value(float64(value))
	case int64:
		return types.Float64Value(float64(value))
	}
	return types.Float64Null()
}

// MapToFrameworkStringSet extracts a string slice from a map and converts it to a Terraform Framework string set.
// Returns empty set for empty arrays, null only if key doesn't exist.
func MapToFrameworkStringSet(ctx context.Context, data map[string]any, key string) types.Set {
	tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Processing key '%s'", key))

	value, exists := data[key]
	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' does not exist in data, returning null", key))
		return types.SetNull(types.StringType)
	}

	tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' exists, value type: %T, value: %v", key, value, value))

	// Handle empty arrays and arrays with values
	if rawSlice, ok := value.([]any); ok {
		tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' is a slice with %d elements", key, len(rawSlice)))

		var strings []string
		for _, item := range rawSlice {
			if str, ok := item.(string); ok {
				strings = append(strings, str)
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' converted to %d strings: %v", key, len(strings), strings))

		// Convert directly to set, always preserving empty arrays as empty sets
		set, diags := types.SetValueFrom(ctx, types.StringType, strings)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' failed to create set, returning empty set. Diagnostics: %v", key, diags))
			// If SetValueFrom fails, return empty set rather than null for existing keys
			emptySet, _ := types.SetValue(types.StringType, []attr.Value{})
			return emptySet
		}

		tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' successfully created set (isEmpty: %v)", key, set.IsNull() || len(set.Elements()) == 0))
		return set
	}

	// Key exists but is not a slice - return null
	tflog.Debug(ctx, fmt.Sprintf("MapToFrameworkStringSet: Key '%s' exists but is not a slice (type: %T), returning null", key, value))
	return types.SetNull(types.StringType)
}

// MapToFrameworkStringList extracts a string slice from a map and converts it to a Terraform Framework string list.
// Returns null if the key doesn't exist or the value is not a []any containing strings.
func MapToFrameworkStringList(ctx context.Context, data map[string]any, key string) types.List {
	if rawSlice, ok := data[key].([]any); ok {
		var strings []string
		for _, item := range rawSlice {
			if str, ok := item.(string); ok {
				strings = append(strings, str)
			}
		}
		// Convert directly to list
		list, diags := types.ListValueFrom(ctx, types.StringType, strings)
		if diags.HasError() {
			return types.ListNull(types.StringType)
		}
		return list
	}
	return types.ListNull(types.StringType)
}
