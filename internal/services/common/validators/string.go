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

// -----------------------------------------------------------------------------------

// illegalCharactersValidator validates that a string does not contain forbidden characters
type illegalCharactersValidator struct {
	forbiddenChars []rune
	description    string
}

// Description describes the validation in plain text formatting.
func (v illegalCharactersValidator) Description(_ context.Context) string {
	if v.description != "" {
		return v.description
	}
	return fmt.Sprintf("string cannot contain the following characters: %s", string(v.forbiddenChars))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v illegalCharactersValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v illegalCharactersValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	for _, char := range value {
		// Check if character is in forbidden list
		for _, forbidden := range v.forbiddenChars {
			if char == forbidden {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Character",
					fmt.Sprintf("String contains forbidden character: %c", char),
				)
				return
			}
		}
	}
}

// IllegalCharactersInString returns a string validator which ensures that the string
// does not contain any of the specified forbidden characters.
func IllegalCharactersInString(forbiddenChars []rune, description string) validator.String {
	return &illegalCharactersValidator{
		forbiddenChars: forbiddenChars,
		description:    description,
	}
}

// asciiOnlyValidator validates that a string contains only ASCII characters (0-127)
type asciiOnlyValidator struct{}

// Description describes the validation in plain text formatting.
func (v asciiOnlyValidator) Description(_ context.Context) string {
	return "string can contain only ASCII characters (0-127)"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v asciiOnlyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v asciiOnlyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	for _, char := range value {
		// Check if character is outside ASCII range 0-127
		if char > 127 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Character",
				fmt.Sprintf("String contains non-ASCII character: %c", char),
			)
			return
		}
	}
}

// ASCIIOnly returns a string validator which ensures that the string
// contains only ASCII characters (0-127).
func ASCIIOnly() validator.String {
	return &asciiOnlyValidator{}
}

// requiredWhenSetContainsValidator validates that a string field is required when a sibling set contains a specific value
type requiredWhenSetContainsValidator struct {
	setFieldName  string
	requiredValue string
}

// Description describes the validation in plain text formatting.
func (v requiredWhenSetContainsValidator) Description(_ context.Context) string {
	return fmt.Sprintf("field is required when %s contains \"%s\"", v.setFieldName, v.requiredValue)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v requiredWhenSetContainsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWhenSetContainsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If value is null or unknown, check if it should be required
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		// Get the sibling set field
		var setField types.Set
		diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName(v.setFieldName), &setField)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		if setField.IsNull() || setField.IsUnknown() {
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

		if hasRequiredValue {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Required Field",
				fmt.Sprintf("field is required when %s contains \"%s\"", v.setFieldName, v.requiredValue),
			)
		}
	}
}

// RequiredWhenSetContains returns a string validator which ensures that the field is required
// when a sibling set field contains a specific value.
func RequiredWhenSetContains(setFieldName, requiredValue string) validator.String {
	return &requiredWhenSetContainsValidator{
		setFieldName:  setFieldName,
		requiredValue: requiredValue,
	}
}

//---------------------------------------------------

// Ensure the implementation satisfies the validator.String interface.
var _ validator.String = requiredWhenEqualsValidator{}

// requiredWhenEqualsValidator is the implementation of the validator.
type requiredWhenEqualsValidator struct {
	dependentField string
	requiredValue  types.String
}

// Description returns a plain-text description of the validator's behavior.
func (v requiredWhenEqualsValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("this attribute is required when %s is set to %s", v.dependentField, v.requiredValue.ValueString())
}

// MarkdownDescription returns a markdown-formatted description of the validator's behavior.
func (v requiredWhenEqualsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWhenEqualsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the attribute being validated is not configured, we don't need to check the dependency.
	// Other validators like `stringvalidator.Required()` will handle if it's required on its own.
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Get the path to the dependent attribute.
	// For example, if we are on "group_id", the parent path is "target", and we want to find "target_type".
	dependentPath := req.Path.ParentPath().AtName(v.dependentField)

	// Get the value of the dependent attribute from the configuration.
	var dependentValue types.String
	diags := req.Config.GetAttribute(ctx, dependentPath, &dependentValue)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If the dependent attribute isn't set to the required value, the validation passes.
	if dependentValue.IsUnknown() || dependentValue.IsNull() || !dependentValue.Equal(v.requiredValue) {
		return
	}

	// At this point, the dependent field has the required value, so this field must be set.
	// Check if the current attribute value is null or an empty string.
	if req.ConfigValue.IsNull() || req.ConfigValue.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Attribute Required",
			fmt.Sprintf("Attribute %q is required because attribute %q is set to %q.", req.Path, dependentPath, v.requiredValue.ValueString()),
		)
	}
}

// RequiredWhenEquals is a factory function that returns a new requiredWhenEqualsValidator.
// It validates that the attribute is not null or empty when another field in the same object
// has a specific string value.
func RequiredWhenEquals(dependentField string, requiredValue types.String) validator.String {
	return requiredWhenEqualsValidator{
		dependentField: dependentField,
		requiredValue:  requiredValue,
	}
}

//---------------------------------------------------

// requiredWithODataValidator validates that OData parameters are only used with filter_type = "odata"
// and that at least one OData parameter is provided when filter_type = "odata"
type requiredWithODataValidator struct {
	odataFieldNames []string
}

// Description describes the validation in plain text formatting.
func (v requiredWithODataValidator) Description(_ context.Context) string {
	return fmt.Sprintf("OData parameters (%s) can only be used when filter_type is 'odata'", strings.Join(v.odataFieldNames, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v requiredWithODataValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v requiredWithODataValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Check if this is the filter_type field and if it's set to "odata"
	if req.Path.String() == "filter_type" && req.ConfigValue.ValueString() == "odata" {
		// Check if at least one OData parameter is provided
		atLeastOneODataParamProvided := false

		for _, odataFieldName := range v.odataFieldNames {
			var odataField basetypes.StringValue
			diags := req.Config.GetAttribute(ctx, path.Root(odataFieldName), &odataField)
			if diags.HasError() {
				continue
			}

			if !odataField.IsNull() && !odataField.IsUnknown() && odataField.ValueString() != "" {
				atLeastOneODataParamProvided = true
				break
			}

			// Check for numeric and boolean OData fields
			if odataFieldName == "odata_top" {
				var numField basetypes.Int64Value
				diags := req.Config.GetAttribute(ctx, path.Root(odataFieldName), &numField)
				if !diags.HasError() && !numField.IsNull() && !numField.IsUnknown() && numField.ValueInt64() > 0 {
					atLeastOneODataParamProvided = true
					break
				}
			}

			if odataFieldName == "odata_count" {
				var boolField basetypes.BoolValue
				diags := req.Config.GetAttribute(ctx, path.Root(odataFieldName), &boolField)
				if !diags.HasError() && !boolField.IsNull() && !boolField.IsUnknown() && boolField.ValueBool() {
					atLeastOneODataParamProvided = true
					break
				}
			}
		}

		if !atLeastOneODataParamProvided {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing OData Parameters",
				fmt.Sprintf("When filter_type is 'odata', at least one of these parameters must be provided: %s", strings.Join(v.odataFieldNames, ", ")),
			)
		}
	}

	// If this is not the filter_type field but an OData field, check if filter_type is "odata"
	for _, odataFieldName := range v.odataFieldNames {
		if req.Path.String() == odataFieldName {
			var filterTypeField basetypes.StringValue
			diags := req.Config.GetAttribute(ctx, path.Root("filter_type"), &filterTypeField)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}

			if !filterTypeField.IsNull() && !filterTypeField.IsUnknown() && filterTypeField.ValueString() != "odata" {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid OData Parameter Usage",
					fmt.Sprintf("OData parameter '%s' can only be used when filter_type is 'odata'", odataFieldName),
				)
			}
		}
	}
}

// ODataParameterValidator returns a string validator which ensures that OData parameters
// are only used with filter_type = "odata" and that at least one OData parameter is provided
// when filter_type = "odata".
func ODataParameterValidator(odataFieldNames ...string) validator.String {
	return &requiredWithODataValidator{
		odataFieldNames: odataFieldNames,
	}
}

//---------------------------------------------------
