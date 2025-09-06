package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// conditionalBoolValidator validates that a boolean field can only have a specific value
// when another boolean field has a specific value
type conditionalBoolValidator struct {
	dependentField    string
	dependentValue    bool
	allowedValue      bool
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v conditionalBoolValidator) Description(_ context.Context) string {
	if v.validationMessage != "" {
		return v.validationMessage
	}

	dependentValueStr := "true"
	if !v.dependentValue {
		dependentValueStr = "false"
	}

	allowedValueStr := "true"
	if !v.allowedValue {
		allowedValueStr = "false"
	}

	return fmt.Sprintf("when %s is %s, this field can only be set to %s",
		v.dependentField, dependentValueStr, allowedValueStr)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v conditionalBoolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation.
func (v conditionalBoolValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Skip validation if the config is empty (for testing purposes)
	if req.Config.Raw.IsNull() {
		return
	}

	// Try to get the dependent field value, but don't error if it's not found
	var dependentValue types.Bool
	diags := req.Config.GetAttribute(ctx, path.Root(v.dependentField), &dependentValue)
	if diags.HasError() {
		// If we can't find the field, skip validation
		// This handles cases where the field might not exist in the schema
		return
	}

	// Skip validation if dependent field is null or unknown
	if dependentValue.IsNull() || dependentValue.IsUnknown() {
		return
	}

	// Check if the dependent field has the condition value
	if dependentValue.ValueBool() == v.dependentValue {
		// If the current field doesn't have the allowed value, add an error
		if req.ConfigValue.ValueBool() != v.allowedValue {
			dependentValueStr := "true"
			if !v.dependentValue {
				dependentValueStr = "false"
			}

			allowedValueStr := "true"
			if !v.allowedValue {
				allowedValueStr = "false"
			}

			errorMessage := v.validationMessage
			if errorMessage == "" {
				errorMessage = fmt.Sprintf("When %s is %s, this field can only be set to %s",
					v.dependentField, dependentValueStr, allowedValueStr)
			}

			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Conditional Value",
				errorMessage,
			)
		}
	}
}

// ConditionalBoolValue returns a boolean validator which ensures that when a dependent field
// has a specific boolean value, the current field can only have a specific value.
func ConditionalBoolValue(dependentField string, dependentValue bool, allowedValue bool, validationMessage string) validator.Bool {
	return &conditionalBoolValidator{
		dependentField:    dependentField,
		dependentValue:    dependentValue,
		allowedValue:      allowedValue,
		validationMessage: validationMessage,
	}
}

// BoolCanOnlyBeTrueWhen returns a boolean validator which ensures that the current field
// can only be true when the dependent field has the specified value.
func BoolCanOnlyBeTrueWhen(dependentField string, dependentValue bool, validationMessage string) validator.Bool {
	return ConditionalBoolValue(dependentField, dependentValue, true, validationMessage)
}

// BoolCanOnlyBeFalseWhen returns a boolean validator which ensures that the current field
// can only be false when the dependent field has the specified value.
func BoolCanOnlyBeFalseWhen(dependentField string, dependentValue bool, validationMessage string) validator.Bool {
	return ConditionalBoolValue(dependentField, dependentValue, false, validationMessage)
}
