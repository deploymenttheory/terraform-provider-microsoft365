package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure that the struct implements the required interface
var _ validator.List = enumListValidator{}

// enumListValidator validates that each item in a list matches one of the expected enum values.
type enumListValidator struct {
	allowedEnumValues []string
}

// Description provides a human-readable description of the validator's purpose.
func (v enumListValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v enumListValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("all list items must be one of: %q", v.allowedEnumValues)
}

// ValidateList performs the validation logic on each item in the list attribute.
func (v enumListValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	// Iterate over each item in the list and validate it
	for i, elem := range request.ConfigValue.Elements() {
		// Ensure the element is a string
		strElem, ok := elem.(types.String)
		if !ok || strElem.IsNull() || strElem.IsUnknown() {
			continue
		}

		itemIsValid := false
		for _, allowedValue := range v.allowedEnumValues {
			if strElem.ValueString() == allowedValue {
				itemIsValid = true
				break
			}
		}

		if !itemIsValid {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				request.Path.AtListIndex(i),
				v.Description(ctx),
				strElem.ValueString(),
			))
		}
	}
}

// EnumValuesList returns an instance of the validator for list values with specified allowed values.
func EnumValuesList(values ...string) validator.List {
	return enumListValidator{
		allowedEnumValues: values,
	}
}
