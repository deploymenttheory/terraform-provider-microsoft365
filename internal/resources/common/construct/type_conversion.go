package construct

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// SetStringProperty sets the value of a string property if the value is not null or unknown.
// It accepts a basetypes.StringValue (Terraform SDK type) and translates it into a pointer to a string for use in the setter function.
func SetStringProperty(value basetypes.StringValue, setter func(*string)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueString()
		setter(&val)
	}
}

// SetBoolProperty sets the value of a bool property if the value is not null or unknown.
// It accepts a basetypes.BoolValue (Terraform SDK type) and translates it into a pointer to a bool for use in the setter function.
func SetBoolProperty(value basetypes.BoolValue, setter func(*bool)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueBool()
		setter(&val)
	}
}

// SetInt32Property sets the value of an int32 property if the value is not null or unknown.
// It accepts a basetypes.Int32Value (Terraform SDK type) and passes it to the setter function.
func SetInt32Property(value basetypes.Int32Value, setter func(*int32)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt32()
		setter(&val)
	}
}

// SetInt64Property sets the value of an int64 property if the value is not null or unknown.
// It accepts a basetypes.Int64Value (Terraform SDK type) and passes it to the setter function.
func SetInt64Property(value basetypes.Int64Value, setter func(*int64)) {
	if !value.IsNull() && !value.IsUnknown() {
		val := value.ValueInt64()
		setter(&val)
	}
}

// ParseEnum parses an enum value and sets it if the value is not null or unknown.
// It accepts a basetypes.StringValue (Terraform SDK type) and uses a parser function to translate the string into an enum type.
// If the value is valid, it casts the parsed value to the expected type T and passes it to the setter function.
func ParseEnum[T any](value basetypes.StringValue, parser func(string) (any, error), setter func(T)) error {
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

// Add this new helper function to your construct package
func SetArrayProperty(values []types.String, setter func([]string)) {
	if len(values) > 0 {
		stringValues := make([]string, 0, len(values))
		for _, v := range values {
			if !v.IsNull() && !v.IsUnknown() {
				stringValues = append(stringValues, v.ValueString())
			}
		}
		if len(stringValues) > 0 {
			setter(stringValues)
		}
	}
}
