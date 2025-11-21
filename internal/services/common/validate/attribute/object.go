package attribute

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// exactlyOneOfMixedTypesValidator validates that exactly one of the specified attributes is set
// This validator supports mixed types (string, int32, bool, set)
type exactlyOneOfMixedTypesValidator struct {
	attributeNames []string
}

// Description returns the validator's description.
func (v exactlyOneOfMixedTypesValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Exactly one of [%s] must be specified", strings.Join(v.attributeNames, ", "))
}

// MarkdownDescription returns the validator's description in Markdown format.
func (v exactlyOneOfMixedTypesValidator) MarkdownDescription(ctx context.Context) string {
	attrNames := make([]string, len(v.attributeNames))
	for i, name := range v.attributeNames {
		attrNames[i] = fmt.Sprintf("`%s`", name)
	}
	return fmt.Sprintf("Exactly one of %s must be specified", strings.Join(attrNames, ", "))
}

// ValidateObject implements the validation logic for object attributes
func (v exactlyOneOfMixedTypesValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	setCount := 0
	var setFields []string

	for _, attrName := range v.attributeNames {
		attrPath := req.Path.AtName(attrName)

		// Try as String
		var stringVal basetypes.StringValue
		if diags := req.Config.GetAttribute(ctx, attrPath, &stringVal); !diags.HasError() {
			if !stringVal.IsNull() && !stringVal.IsUnknown() {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}

		// Try as Int32
		var int32Val basetypes.Int32Value
		if diags := req.Config.GetAttribute(ctx, attrPath, &int32Val); !diags.HasError() {
			if !int32Val.IsNull() && !int32Val.IsUnknown() {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}

		// Try as Int64
		var int64Val basetypes.Int64Value
		if diags := req.Config.GetAttribute(ctx, attrPath, &int64Val); !diags.HasError() {
			if !int64Val.IsNull() && !int64Val.IsUnknown() {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}

		// Try as Bool
		var boolVal basetypes.BoolValue
		if diags := req.Config.GetAttribute(ctx, attrPath, &boolVal); !diags.HasError() {
			if !boolVal.IsNull() && !boolVal.IsUnknown() {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}

		// Try as Set
		var setVal basetypes.SetValue
		if diags := req.Config.GetAttribute(ctx, attrPath, &setVal); !diags.HasError() {
			if !setVal.IsNull() && !setVal.IsUnknown() && len(setVal.Elements()) > 0 {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}

		// Try as List
		var listVal basetypes.ListValue
		if diags := req.Config.GetAttribute(ctx, attrPath, &listVal); !diags.HasError() {
			if !listVal.IsNull() && !listVal.IsUnknown() && len(listVal.Elements()) > 0 {
				setCount++
				setFields = append(setFields, attrName)
				continue
			}
		}
	}

	if setCount == 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Attribute",
			fmt.Sprintf("Exactly one of these attributes must be specified: %s",
				strings.Join(v.attributeNames, ", ")),
		)
		return
	}

	if setCount > 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Conflicting Attributes",
			fmt.Sprintf("Only one of these attributes can be specified: %s. Found multiple: %s",
				strings.Join(v.attributeNames, ", "),
				strings.Join(setFields, ", ")),
		)
	}
}

// ExactlyOneOfMixedTypes returns a validator that ensures exactly one of the specified attributes is set.
// This validator supports attributes of different types (string, int32, int64, bool, set, list).
//
// Example usage in a schema:
//
//	"my_object": schema.SingleNestedAttribute{
//	    Validators: []validator.Object{
//	        attribute.ExactlyOneOfMixedTypes("string_value", "int_value", "bool_value"),
//	    },
//	}
func ExactlyOneOfMixedTypes(attributeNames ...string) validator.Object {
	return &exactlyOneOfMixedTypesValidator{
		attributeNames: attributeNames,
	}
}
