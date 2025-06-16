package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringSetValidator validates that a string set only contains allowed values.
type stringSetValidator struct {
	allowedValues []string
}

// Description describes the validation in plain text formatting.
func (v stringSetValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %v", v.allowedValues)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v stringSetValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateSet performs the validation.
func (v stringSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	setValues := make(map[string]struct{})
	for _, allowed := range v.allowedValues {
		setValues[allowed] = struct{}{}
	}

	elements := req.ConfigValue.Elements()
	for _, element := range elements {
		str, ok := element.(types.String)
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid Set Element",
				"Set element is not a string type",
			)
			return
		}

		if str.IsNull() || str.IsUnknown() {
			continue
		}

		if _, ok := setValues[str.ValueString()]; !ok {
			resp.Diagnostics.AddError(
				"Invalid Set Element Value",
				fmt.Sprintf("Set element value must be one of: %v", v.allowedValues),
			)
			return
		}
	}
}

// StringSetAllowedValues returns a Set validator which ensures that any configured
// string set value matches one of the allowed values exactly.
func StringSetAllowedValues(allowedValues ...string) validator.Set {
	return &stringSetValidator{
		allowedValues: allowedValues,
	}
}

// blockRequiresSetValueValidator validates that a block can only exist if a specific value exists in a sibling set field
type blockRequiresSetValueValidator struct {
	setFieldName   string
	requiredValue  string
	blockFieldName string
}

// Description describes the validation in plain text formatting.
func (v blockRequiresSetValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("block can only be specified when \"%s\" is included in %s", v.requiredValue, v.setFieldName)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v blockRequiresSetValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateObject performs the validation.
func (v blockRequiresSetValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get the sibling set field
	var setField types.Set
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName(v.setFieldName), &setField)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if setField.IsNull() || setField.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Configuration",
			fmt.Sprintf("%s block cannot be specified when %s is not configured", v.blockFieldName, v.setFieldName),
		)
		return
	}

	// Check if the required value exists in the set
	hasRequiredValue := false
	elements := setField.Elements()
	for _, element := range elements {
		str, ok := element.(types.String)
		if !ok {
			continue
		}
		if !str.IsNull() && !str.IsUnknown() && str.ValueString() == v.requiredValue {
			hasRequiredValue = true
			break
		}
	}

	if !hasRequiredValue {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Configuration",
			fmt.Sprintf("%s block is specified but \"%s\" is not included in %s. Add \"%s\" to %s or remove the %s block.",
				v.blockFieldName, v.requiredValue, v.setFieldName, v.requiredValue, v.setFieldName, v.blockFieldName),
		)
	}
}

// BlockRequiresSetValue returns an Object validator which ensures that
// a block can only be specified when a specific value exists in a sibling set field.
func BlockRequiresSetValue(setFieldName, requiredValue, blockFieldName string) validator.Object {
	return &blockRequiresSetValueValidator{
		setFieldName:   setFieldName,
		requiredValue:  requiredValue,
		blockFieldName: blockFieldName,
	}
}
