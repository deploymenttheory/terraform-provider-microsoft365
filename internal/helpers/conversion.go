package helpers

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
	return types.StringValue(t.Format(TimeFormatRFC3339))
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
