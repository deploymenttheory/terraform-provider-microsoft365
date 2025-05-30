package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// requiredWithValidator validates that a string field is set when another field has a specific value
type requiredWithValidator struct {
	fieldName  string
	fieldValue string
}

// Description describes the validation in plain text formatting.
func (v requiredWithValidator) Description(_ context.Context) string {
	return fmt.Sprintf("field is required when %s is %s", v.fieldName, v.fieldValue)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v requiredWithValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWithValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If value is being reset to null/empty, check the condition field
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		var conditionField types.String
		diags := req.Config.GetAttribute(ctx, path.Root(v.fieldName), &conditionField)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		if conditionField.ValueString() == v.fieldValue {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Required Field",
				fmt.Sprintf("field is required when %s is %s", v.fieldName, v.fieldValue),
			)
		}
	}
}

// RequiredWith returns a string validator which ensures that the field is set
// when another field matches a specific value.
func RequiredWith(fieldName string, fieldValue string) validator.String {
	return &requiredWithValidator{
		fieldName:  fieldName,
		fieldValue: fieldValue,
	}
}

//---------------------------------------------------

// MutuallyExclusiveAttributesValidator checks if only one of the specified attributes is set
type MutuallyExclusiveAttributesValidator struct {
	AttributeNames []string
}

// Description returns the validator's description.
func (v MutuallyExclusiveAttributesValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures that only one of the attributes [%s] is set", strings.Join(v.AttributeNames, ", "))
}

// MarkdownDescription returns the validator's description in Markdown format.
func (v MutuallyExclusiveAttributesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateObject implements validator logic
func (v MutuallyExclusiveAttributesValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If less than 2 attributes to check, validation is unnecessary
	if len(v.AttributeNames) < 2 {
		return
	}

	// Count attributes that are set (non-empty strings)
	setCount := 0
	var setFields []string

	for _, attrName := range v.AttributeNames {
		// Use simple individual string attribute checks
		var value basetypes.StringValue

		// Create a proper path for the attribute
		attrPath := req.Path.AtName(attrName)

		diags := req.Config.GetAttribute(ctx, attrPath, &value)

		// Skip attributes that don't exist or can't be accessed
		if diags.HasError() {
			continue
		}

		// Check if attribute is set (not null and not empty)
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			setCount++
			setFields = append(setFields, attrName)
		}
	}

	// If more than one attribute is set, add an error
	if setCount > 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Mutually Exclusive Attributes",
			fmt.Sprintf("Only one of these attributes can be specified: %s. Found multiple: %s",
				strings.Join(v.AttributeNames, ", "),
				strings.Join(setFields, ", ")),
		)
	}
}

// ExactlyOneOf returns a validator that ensures exactly one of the specified attributes is set.
// This validator works on nested attributes within a block.
func ExactlyOneOf(attributeNames ...string) validator.Object {
	return &MutuallyExclusiveAttributesValidator{
		AttributeNames: attributeNames,
	}
}

// -----------------------------------------------------------------------------------

// stringLengthValidator validates that a string field does not exceed a maximum length
type stringLengthValidator struct {
	maxLength int
}

// Description describes the validation in plain text formatting.
func (v stringLengthValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must not exceed %d characters", v.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v stringLengthValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v stringLengthValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) > v.maxLength {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"String Too Long",
			fmt.Sprintf("String length %d exceeds maximum allowed length of %d characters", len(value), v.maxLength),
		)
	}
}

// StringLengthAtMost returns a string validator which ensures that the string
// does not exceed the specified maximum length.
func StringLengthAtMost(maxLength int) validator.String {
	return &stringLengthValidator{
		maxLength: maxLength,
	}
}
