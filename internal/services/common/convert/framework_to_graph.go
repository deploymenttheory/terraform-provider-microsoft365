package convert

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// ============================================================================
// TERRAFORM FRAMEWORK → GRAPH SDK (Write Operations)
// Used in Create/Update() methods when sending Terraform config to Graph API
// ============================================================================

// FrameworkToGraphString sets a Graph SDK string property from a Terraform Framework string.
// Only sets the value if it's not null or unknown.
func FrameworkToGraphString(value basetypes.StringValue, setter func(*string)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueString()
		setter(&val)
	}
}

// FrameworkToGraphBool sets a Graph SDK bool property from a Terraform Framework bool.
// Only sets the value if it's not null or unknown.
func FrameworkToGraphBool(value basetypes.BoolValue, setter func(*bool)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueBool()
		setter(&val)
	}
}

// FrameworkToGraphInt32 sets a Graph SDK int32 property from a Terraform Framework int32.
// Only sets the value if it's not null or unknown.
func FrameworkToGraphInt32(value basetypes.Int32Value, setter func(*int32)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt32()
		setter(&val)
	}
}

// FrameworkToGraphInt64 sets a Graph SDK int64 property from a Terraform Framework int64.
// Only sets the value if it's not null or unknown.
func FrameworkToGraphInt64(value basetypes.Int64Value, setter func(*int64)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt64()
		setter(&val)
	}
}

// FrameworkToGraphTime parses a Terraform Framework string as RFC3339 time and sets a Graph SDK time property.
// Returns an error if parsing fails. No-op if the value is null, unknown, or empty.
func FrameworkToGraphTime(value basetypes.StringValue, setter func(*time.Time)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	dateStr := value.ValueString()
	if dateStr == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse time string: %s", err)
	}

	setter(&parsed)
	return nil
}

// FrameworkToGraphDateOnly parses a Terraform Framework string as a date and sets a Graph SDK DateOnly property.
// Returns an error if parsing fails. No-op if the value is null, unknown, or empty.
func FrameworkToGraphDateOnly(value basetypes.StringValue, setter func(*serialization.DateOnly)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	dateStr := value.ValueString()
	if dateStr == "" {
		return nil
	}

	parsedDate, err := serialization.ParseDateOnly(dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse date string: %s", err)
	}

	setter(parsedDate)
	return nil
}

// Supports various time formats with different nanosecond precision levels (HH:MM:SS, HH:MM:SS.fff, etc.).
// Returns an error if parsing fails. No-op if the value is null, unknown, or empty.
func FrameworkToGraphTimeOnly(value types.String, setter func(*serialization.TimeOnly)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	timeStr := strings.TrimSpace(value.ValueString())
	if timeStr == "" {
		return nil
	}

	// Handle HH:MM format cases by converting to HH:MM:SS
	if matched, _ := regexp.MatchString(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`, timeStr); matched {
		timeStr = timeStr + ":00"
	}

	timeOnly, _, err := serialization.ParseTimeOnlyWithPrecision(timeStr)
	if err != nil {
		return fmt.Errorf("failed to parse time string '%s': expected format HH:MM or HH:MM:SS (e.g., '14:30' or '14:30:00'), got error: %v", value.ValueString(), err)
	}

	if timeOnly != nil {
		setter(timeOnly)
	} else {
		return fmt.Errorf("parsed time resulted in nil TimeOnly object for input '%s'", value.ValueString())
	}

	return nil
}

// FrameworkToGraphTimeOnlyWithPrecision parses a Terraform Framework string as time with explicit precision control.
// Supports various time formats and allows specifying the desired output precision.
// precision: 0-9, where 0 = HH:MM:SS, 1 = HH:MM:SS.f, 2 = HH:MM:SS.ff, etc.
func FrameworkToGraphTimeOnlyWithPrecision(value types.String, precision int, setter func(*serialization.TimeOnly)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	timeStr := value.ValueString()
	if timeStr == "" {
		return nil
	}

	timeOnly, detectedPrecision, err := serialization.ParseTimeOnlyWithPrecision(timeStr)
	if err != nil {
		return fmt.Errorf("failed to parse time string '%s': %v", timeStr, err)
	}

	if timeOnly != nil {
		// If we need a different precision than detected, we can create a new TimeOnly
		// with the desired precision by formatting and re-parsing
		if precision != detectedPrecision && precision >= 0 && precision <= 9 {
			// Format with desired precision and re-parse
			formattedTime := timeOnly.StringWithPrecision(precision)
			timeOnly, _, err = serialization.ParseTimeOnlyWithPrecision(formattedTime)
			if err != nil {
				return fmt.Errorf("failed to reformat time with precision %d: %v", precision, err)
			}
		}
		setter(timeOnly)
	}

	return nil
}

// FrameworkToGraphISODuration parses a Terraform Framework string as ISO 8601 duration and sets a Graph SDK ISODuration property.
// Returns an error if parsing fails. No-op if the value is null, unknown, or empty.
func FrameworkToGraphISODuration(value basetypes.StringValue, setter func(*serialization.ISODuration)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	isoDuration, err := serialization.ParseISODuration(value.ValueString())
	if err != nil {
		return fmt.Errorf("error parsing ISO 8601 duration: %v", err)
	}
	setter(isoDuration)
	return nil
}

// FrameworkToGraphUUID parses a Terraform Framework string as UUID and sets a Graph SDK UUID property.
// Returns an error if parsing fails. No-op if the value is null, unknown, or empty.
func FrameworkToGraphUUID(value basetypes.StringValue, setter func(*uuid.UUID)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	raw := value.ValueString()
	if raw == "" {
		return nil
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid UUID: %s", err)
	}

	setter(&parsed)
	return nil
}

// FrameworkToGraphBytes converts a Terraform Framework string to bytes and sets a Graph SDK byte slice property.
// Only sets the value if it's not null or unknown.
func FrameworkToGraphBytes(value basetypes.StringValue, setter func([]byte)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := []byte(value.ValueString())
		setter(val)
	}
}

// FrameworkToGraphEnum parses a Terraform Framework string as an enum and sets a Graph SDK enum property.
// Uses a parser function to convert the string to the enum type.
// Returns an error if parsing or type assertion fails. No-op if the value is null or unknown.
func FrameworkToGraphEnum[T any](value basetypes.StringValue, parser func(string) (any, error), setter func(T)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	enumVal, err := parser(value.ValueString())
	if err != nil {
		return fmt.Errorf("failed to parse enum: %v", err)
	}

	typedEnumVal, ok := enumVal.(T)
	if !ok {
		return fmt.Errorf("failed to cast parsed value to expected type %T", enumVal)
	}

	setter(typedEnumVal)
	return nil
}

// FrameworkToGraphBitmaskEnum parses a Terraform Framework string as a bitmask-style enum and sets a Graph SDK enum property.
// Expects the parser to return a pointer to the enum type.
// Returns an error if parsing or type assertion fails. No-op if the value is null or unknown.
func FrameworkToGraphBitmaskEnum[T any](value basetypes.StringValue, parser func(string) (any, error), setter func(*T)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	enumVal, err := parser(value.ValueString())
	if err != nil {
		return fmt.Errorf("failed to parse enum: %v", err)
	}
	if enumVal == nil {
		return nil // silently ignore like Microsoft's parser
	}

	typed, ok := enumVal.(*T)
	if !ok {
		return fmt.Errorf("failed to cast parsed value to expected type")
	}

	setter(typed)
	return nil
}

// FrameworkToGraphStringList converts a Terraform Framework list to a string slice and sets a Graph SDK string slice property.
// Returns an error if list elements are not strings. Sets nil if the list is null or unknown.
func FrameworkToGraphStringList(ctx context.Context, list types.List, setter func([]string)) error {
	if list.IsNull() || list.IsUnknown() {
		setter(nil)
		return nil
	}

	elements := list.Elements()
	result := make([]string, 0, len(elements))
	for i, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			return fmt.Errorf("unexpected element type at index %d: %T", i, elem)
		}

		if !strVal.IsNull() && !strVal.IsUnknown() {
			result = append(result, strVal.ValueString())
		}
	}

	setter(result)
	return nil
}

// FrameworkToGraphStringSet converts a Terraform Framework set to a string slice and sets a Graph SDK string slice property.
// Returns an error if set elements are not strings. Sets nil if the set is null or unknown.
func FrameworkToGraphStringSet(ctx context.Context, set types.Set, setter func([]string)) error {
	if set.IsNull() || set.IsUnknown() {
		setter(nil)
		return nil
	}

	elements := set.Elements()
	result := make([]string, 0, len(elements))
	for i, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			return fmt.Errorf("unexpected element type at index %d: %T", i, elem)
		}

		if !strVal.IsNull() && !strVal.IsUnknown() {
			result = append(result, strVal.ValueString())
		}
	}

	setter(result)
	return nil
}

// FrameworkToGraphObjectsFromStringSet is a generic function that converts a Terraform Framework string set to objects.
// Extracts string values from the set, passes them to a converter function to transform them into
// the desired object type, and then sets them using the provided setter function.
func FrameworkToGraphObjectsFromStringSet[T any](
	ctx context.Context,
	set types.Set,
	converter func(context.Context, []string) []T,
	setter func([]T)) error {

	if set.IsNull() || set.IsUnknown() {
		setter(nil)
		return nil
	}

	var stringValues []string
	diags := set.ElementsAs(ctx, &stringValues, false)
	if diags.HasError() {
		return fmt.Errorf("failed to extract string values: %s", diags.Errors())
	}

	objects := converter(ctx, stringValues)
	setter(objects)
	return nil
}

// FrameworkToGraphBitmaskEnumFromSet converts a Terraform Framework set of strings to a bitmask enum.
// This is useful for APIs that use bitmask enums with String() methods that return comma-separated values.
// The function joins the set elements with commas, parses the result using the provided parser function,
// and sets the resulting enum using the provided setter function.
// Returns an error if parsing fails. No-op if the set is null or unknown.
func FrameworkToGraphBitmaskEnumFromSet[T any](
	ctx context.Context,
	set types.Set,
	parser func(string) (any, error),
	setter func(*T)) error {

	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	// Extract string values from the set
	elements := set.Elements()
	if len(elements) == 0 {
		return nil
	}

	// Convert to string slice
	var stringValues []string
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
			stringValues = append(stringValues, strVal.ValueString())
		}
	}

	if len(stringValues) == 0 {
		return nil
	}

	// Join with commas and parse
	joinedStr := strings.Join(stringValues, ",")
	result, err := parser(joinedStr)
	if err != nil {
		return fmt.Errorf("failed to parse bitmask enum: %v", err)
	}

	// Type assert and set
	if result == nil {
		return nil
	}

	typed, ok := result.(*T)
	if !ok {
		return fmt.Errorf("failed to cast parsed value to expected type")
	}

	setter(typed)
	return nil
}
