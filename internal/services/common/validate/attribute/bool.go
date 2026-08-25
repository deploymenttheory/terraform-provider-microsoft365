package attribute

import (
	"context"
	"fmt"
	"slices"
	"strings"

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

// conditionalBoolPresenceValidator validates that a boolean field may only be configured at
// all - with either value - depending on the value held by another string attribute.
//
// This differs from conditionalStringBoolValidator, which constrains which *value* the
// boolean may take. Here the constraint is on presence: some Graph settings are rejected
// outright by the service unless a sibling discriminator holds a particular value, so the
// provider must be able to leave the field out of the request entirely.
type conditionalBoolPresenceValidator struct {
	dependentPath     path.Path
	values            []string
	allowList         bool
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v conditionalBoolPresenceValidator) Description(_ context.Context) string {
	if v.validationMessage != "" {
		return v.validationMessage
	}

	if v.allowList {
		return fmt.Sprintf("this field can only be set when %s is one of: %s",
			v.dependentPath, strings.Join(v.values, ", "))
	}

	return fmt.Sprintf("this field cannot be set when %s is one of: %s",
		v.dependentPath, strings.Join(v.values, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v conditionalBoolPresenceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation.
func (v conditionalBoolPresenceValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// An absent or not yet known field cannot violate a presence constraint.
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
			errorMessage = fmt.Sprintf("This field cannot be set when %s is %q. It is only supported when %s is one of: %s.",
				v.dependentPath, dependentValue.ValueString(), v.dependentPath, strings.Join(v.values, ", "))
		} else {
			errorMessage = fmt.Sprintf("This field cannot be set when %s is %q.",
				v.dependentPath, dependentValue.ValueString())
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Attribute Combination",
		errorMessage,
	)
}

// BoolSupportedOnlyWhenStringIn returns a boolean validator which ensures that the field is
// only configured when the string attribute at dependentPath holds one of the given values.
// Setting the field to either true or false is an error for any other value.
//
// Use it for Graph settings the service accepts only in a specific mode, where the provider
// must omit the field from the request otherwise. Such attributes must not carry a schema
// Default, since a default is materialised into the plan for every configuration and would
// then be sent in every request.
//
// Example usage - is_removable is only supported for a required install intent:
//
//	validate.BoolSupportedOnlyWhenStringIn(
//	    path.Root("intent"),
//	    []string{"required"},
//	    "is_removable is only supported when intent is 'required'.",
//	)
func BoolSupportedOnlyWhenStringIn(dependentPath path.Path, values []string, validationMessage string) validator.Bool {
	return &conditionalBoolPresenceValidator{
		dependentPath:     dependentPath,
		values:            values,
		allowList:         true,
		validationMessage: validationMessage,
	}
}

// BoolUnsupportedWhenStringIn returns a boolean validator which ensures that the field is not
// configured when the string attribute at dependentPath holds one of the given values. It is
// the inverse of BoolSupportedOnlyWhenStringIn, for the cases where the unsupported values
// are the shorter and more stable list.
//
// Example usage - uninstall_on_device_removal is rejected for an uninstall intent:
//
//	validate.BoolUnsupportedWhenStringIn(
//	    path.Root("intent"),
//	    []string{"uninstall"},
//	    "uninstall_on_device_removal is not supported when intent is 'uninstall'.",
//	)
func BoolUnsupportedWhenStringIn(dependentPath path.Path, values []string, validationMessage string) validator.Bool {
	return &conditionalBoolPresenceValidator{
		dependentPath:     dependentPath,
		values:            values,
		allowList:         false,
		validationMessage: validationMessage,
	}
}
