package constructors

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// SetStringProperty sets the value of a string property if the value is not null or unknown.
// It accepts a basetypes.StringValue (Terraform SDK type) and translates it into a pointer
// to a string for use in the msgraph-sdk-go setter function.
func SetStringProperty(value basetypes.StringValue, setter func(*string)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueString()
		setter(&val)
	}
}

// SetBoolProperty sets the value of a bool property if the value is not null or unknown.
// It accepts a basetypes.BoolValue (Terraform SDK type) and translates it into a pointer
// to a bool for use in the setter function.
func SetBoolProperty(value basetypes.BoolValue, setter func(*bool)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueBool()
		setter(&val)
	}
}

// SetInt32Property sets the value of an int32 property if the value is not null or unknown.
// It accepts a basetypes.Int32Value (Terraform SDK type) and passes it to the msgraph-sdk-go
// setter function.
func SetInt32Property(value basetypes.Int32Value, setter func(*int32)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt32()
		setter(&val)
	}
}

// SetInt64Property sets the value of an int64 property if the value is not null or unknown.
// It accepts a basetypes.Int64Value (Terraform SDK type) and passes it to the msgraph-sdk-go
// setter function.
func SetInt64Property(value basetypes.Int64Value, setter func(*int64)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt64()
		setter(&val)
	}
}

// SetEnumProperty parses an enum value and sets it if the value is not null or unknown.
// It accepts a basetypes.StringValue (Terraform SDK type) and uses a parser function to
// translate the string into an enum type. If the value is valid, it casts the parsed value
// to the expected type T and passes it to the msgraph-sdk-go setter function.
func SetEnumProperty[T any](value basetypes.StringValue, parser func(string) (any, error), setter func(T)) error {
	if !value.IsNull() && !value.IsUnknown() {

		enumVal, err := parser(value.ValueString())
		if err != nil {
			return fmt.Errorf("failed to parse enum: %v", err)
		}

		// Perform the type assertion to convert from `any` to the expected type `T`
		typedEnumVal, ok := enumVal.(T)
		if !ok {
			return fmt.Errorf("failed to cast parsed value to expected type %T", enumVal)
		}

		setter(typedEnumVal)
	}
	return nil
}

// SetStringList constructs and sets a slice of strings from a Terraform ListAttribute.
// It handles null or unknown values and converts each element to a string and passes it to
// the msgraph-sdk-go setter function.
func SetStringList(ctx context.Context, list types.List, setter func([]string)) error {
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

// SetStringSet constructs and sets a slice of strings from a Terraform SetAttribute.
// It handles null or unknown values and converts each element to a string and passes it to
// the msgraph-sdk-go setter function.
func SetStringSet(ctx context.Context, set types.Set, setter func([]string)) error {
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

// SetBytesProperty sets the value of a byte slice property if the value is not null or unknown.
// It converts a basetypes.StringValue (Terraform SDK type) to a []byte and passes it to the setter function.
func SetBytesProperty(value basetypes.StringValue, setter func([]byte)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := []byte(value.ValueString())
		setter(val)
	}
}

// SetISODurationProperty parses an ISO 8601 duration string and sets the value if valid.
// It accepts a basetypes.StringValue (Terraform SDK type), parses it into ISODuration, and passes it to the setter function.
func SetISODurationProperty(value basetypes.StringValue, setter func(*serialization.ISODuration)) error {
	if !value.IsNull() && !value.IsUnknown() {
		isoDuration, err := serialization.ParseISODuration(value.ValueString())
		if err != nil {
			return fmt.Errorf("error parsing ISO 8601 duration: %v", err)
		}
		setter(isoDuration)
	}
	return nil
}

// SetObjectsFromStringSet is a generic function that constructs objects from a Terraform SetAttribute.
// It extracts string values from the set, passes them to a converter function to transform them into
// the desired object type, and then sets them using the provided setter function.
func SetObjectsFromStringSet[T any](
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

// StringToTime parses a string value into a time.Time if the value is not null or unknown,
// and sets it using the provided setter function.
func StringToTime(value basetypes.StringValue, setter func(*time.Time)) error {
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

// StringToTimeOnly converts a string in HH:MM:SS[.mmmmmmm] format to a TimeOnly type for the Microsoft Graph SDK.
// It handles null or unknown values by returning nil, which is appropriate for optional time fields.
// The function accepts a basetypes.StringValue (Terraform SDK type) and returns a *serialization.TimeOnly.
func StringToTimeOnly(value types.String, setter func(*serialization.TimeOnly)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	timeStr := value.ValueString()
	if timeStr == "" {
		return nil
	}

	// Parse the time string to a TimeOnly object
	timeOnly, err := serialization.ParseTimeOnly(timeStr)
	if err != nil {
		return fmt.Errorf("failed to parse time string '%s': %v", timeStr, err)
	}

	// Set the value using the provided setter function
	setter(timeOnly)
	return nil
}
