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

// conditionalStringBoolValidator validates that a boolean field can only have a specific value
// when another string field has a specific value
type conditionalStringBoolValidator struct {
	dependentField    string
	dependentValue    string
	allowedValue      bool
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v conditionalStringBoolValidator) Description(_ context.Context) string {
	if v.validationMessage != "" {
		return v.validationMessage
	}

	allowedValueStr := "true"
	if !v.allowedValue {
		allowedValueStr = "false"
	}

	return fmt.Sprintf("when %s is %s, this field can only be set to %s",
		v.dependentField, v.dependentValue, allowedValueStr)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v conditionalStringBoolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation.
func (v conditionalStringBoolValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Skip validation if the config is empty (for testing purposes)
	if req.Config.Raw.IsNull() {
		return
	}

	// Try to get the dependent field value, but don't error if it's not found
	var dependentValue types.String
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
	if dependentValue.ValueString() == v.dependentValue {
		// If the current field doesn't have the allowed value, add an error
		if req.ConfigValue.ValueBool() != v.allowedValue {
			errorMessage := v.validationMessage
			if errorMessage == "" {
				allowedValueStr := "true"
				if !v.allowedValue {
					allowedValueStr = "false"
				}
				errorMessage = fmt.Sprintf("When %s is %s, this field can only be set to %s",
					v.dependentField, v.dependentValue, allowedValueStr)
			}

			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Conditional Value",
				errorMessage,
			)
		}
	}
}

// ConditionalStringBoolValue returns a boolean validator which ensures that when a dependent string field
// has a specific string value, the current field can only have a specific boolean value.
func ConditionalStringBoolValue(dependentField string, dependentValue string, allowedValue bool, validationMessage string) validator.Bool {
	return &conditionalStringBoolValidator{
		dependentField:    dependentField,
		dependentValue:    dependentValue,
		allowedValue:      allowedValue,
		validationMessage: validationMessage,
	}
}

// BoolCanOnlyBeTrueWhenStringEquals returns a boolean validator which ensures that the current field
// can only be true when the dependent string field has the specified value.
func BoolCanOnlyBeTrueWhenStringEquals(dependentField string, dependentValue string, validationMessage string) validator.Bool {
	return ConditionalStringBoolValue(dependentField, dependentValue, true, validationMessage)
}

// BoolCanOnlyBeFalseWhenStringEquals returns a boolean validator which ensures that the current field
// can only be false when the dependent string field has the specified value.
func BoolCanOnlyBeFalseWhenStringEquals(dependentField string, dependentValue string, validationMessage string) validator.Bool {
	return ConditionalStringBoolValue(dependentField, dependentValue, false, validationMessage)
}

// mutuallyExclusiveBoolValidator validates that two boolean fields cannot both be true at the same time
type mutuallyExclusiveBoolValidator struct {
	otherField        string
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v mutuallyExclusiveBoolValidator) Description(_ context.Context) string {
	if v.validationMessage != "" {
		return v.validationMessage
	}

	return fmt.Sprintf("this field and %s cannot both be set to true", v.otherField)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v mutuallyExclusiveBoolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation.
func (v mutuallyExclusiveBoolValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// Skip validation if the current value is null, unknown, or false
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || !req.ConfigValue.ValueBool() {
		return
	}

	// Skip validation if the config is empty (for testing purposes)
	if req.Config.Raw.IsNull() {
		return
	}

	// Try to get the other field value
	var otherValue types.Bool
	diags := req.Config.GetAttribute(ctx, path.Root(v.otherField), &otherValue)
	if diags.HasError() {
		// If we can't find the field, skip validation
		return
	}

	// Skip validation if other field is null or unknown
	if otherValue.IsNull() || otherValue.IsUnknown() {
		return
	}

	// If both fields are true, add an error
	if otherValue.ValueBool() {
		errorMessage := v.validationMessage
		if errorMessage == "" {
			errorMessage = fmt.Sprintf("The fields cannot both be set to true. Either this field or %s must be false.", v.otherField)
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Mutually Exclusive Fields",
			errorMessage,
		)
	}
}

// MutuallyExclusiveBool returns a boolean validator which ensures that the current field
// and another boolean field cannot both be true at the same time.
func MutuallyExclusiveBool(otherField string, validationMessage string) validator.Bool {
	return &mutuallyExclusiveBoolValidator{
		otherField:        otherField,
		validationMessage: validationMessage,
	}
}
