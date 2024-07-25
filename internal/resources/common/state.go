package common

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// SetStringValueFromAttributes sets a string value from the given attribute map if the key exists and is not null.
// It takes a map of attributes, a key to look for, and a setter function that sets the value if found.
func SetStringValueFromAttributes(attrs map[string]attr.Value, key string, setter func(*string)) {
	if v, ok := attrs[key].(types.String); ok && !v.IsNull() {
		str := v.ValueString()
		setter(&str)
	}
}

// SetParsedValueFromAttributes sets a parsed value from the given attribute map if the key exists and is not null.
// It takes a map of attributes, a key to look for, a setter function to set the parsed value, and a parser function
// to convert the string value to the desired type. It returns an error if parsing fails.
func SetParsedValueFromAttributes[T any](attrs map[string]attr.Value, key string, setter func(*T), parser func(string) (interface{}, error)) error {
	if v, ok := attrs[key].(types.String); ok && !v.IsNull() {
		str := v.ValueString()
		parsedValue, err := parser(str)
		if err != nil {
			return err
		}
		if parsedValue != nil {
			setter(parsedValue.(*T))
		}
	}
	return nil
}

// safeDeref safely dereferences a string pointer.
// It returns an empty string if the pointer is nil,
// otherwise it returns the dereferenced string value.
func SafeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// safeEnumString safely converts an enum to its string representation.
// It returns an empty string if the enum is nil,
// otherwise it calls the String() method on the enum.
// This function expects the input to implement the fmt.Stringer interface.
func SafeEnumString(e fmt.Stringer) string {
	if e == nil {
		return ""
	}
	return e.String()
}
