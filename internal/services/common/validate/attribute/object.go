package attribute

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// conditionalObjectPresenceValidator validates that a nested object may only be configured at
// all depending on the value held by another string attribute.
//
// The object counterpart of conditionalBoolPresenceValidator: some Graph settings are rejected
// outright by the service unless a sibling discriminator holds a particular value, so the
// provider must be able to leave the block out of the request entirely.
type conditionalObjectPresenceValidator struct {
	dependentPath     path.Path
	values            []string
	allowList         bool
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v conditionalObjectPresenceValidator) Description(_ context.Context) string {
	if v.validationMessage != "" {
		return v.validationMessage
	}

	if v.allowList {
		return fmt.Sprintf("this block can only be set when %s is one of: %s",
			v.dependentPath, strings.Join(v.values, ", "))
	}

	return fmt.Sprintf("this block cannot be set when %s is one of: %s",
		v.dependentPath, strings.Join(v.values, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v conditionalObjectPresenceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateObject performs the validation.
func (v conditionalObjectPresenceValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// An absent or not yet known block cannot violate a presence constraint.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if req.Config.Raw.IsNull() {
		return
	}

	var dependentValue types.String
	if diags := req.Config.GetAttribute(ctx, v.dependentPath, &dependentValue); diags.HasError() {
		// The dependent attribute is not part of this schema; nothing to validate against.
		return
	}

	// The constraint cannot be evaluated until the dependent value is known. The service
	// still rejects an unsupported combination at apply time, with its own error.
	if dependentValue.IsNull() || dependentValue.IsUnknown() {
		return
	}

	matched := slices.Contains(v.values, dependentValue.ValueString())
	if matched == v.allowList {
		return
	}

	errorMessage := v.validationMessage
	if errorMessage == "" {
		if v.allowList {
			errorMessage = fmt.Sprintf("This block cannot be set when %s is %q. It is only supported when %s is one of: %s.",
				v.dependentPath, dependentValue.ValueString(), v.dependentPath, strings.Join(v.values, ", "))
		} else {
			errorMessage = fmt.Sprintf("This block cannot be set when %s is %q.",
				v.dependentPath, dependentValue.ValueString())
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Attribute Combination",
		errorMessage,
	)
}

// ObjectSupportedOnlyWhenStringIn returns an object validator which ensures that the block is
// only configured when the string attribute at dependentPath holds one of the given values.
//
// Use it for Graph settings blocks the service accepts only in a specific mode, where the
// provider must omit the block from the request otherwise.
//
// Example usage - auto_update_settings is only supported for an available install intent:
//
//	validate.ObjectSupportedOnlyWhenStringIn(
//	    path.Root("intent"),
//	    []string{"available"},
//	    "auto_update_settings is only supported when intent is 'available'.",
//	)
func ObjectSupportedOnlyWhenStringIn(dependentPath path.Path, values []string, validationMessage string) validator.Object {
	return &conditionalObjectPresenceValidator{
		dependentPath:     dependentPath,
		values:            values,
		allowList:         true,
		validationMessage: validationMessage,
	}
}
