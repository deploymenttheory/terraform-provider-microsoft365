package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure that the struct implements the required interface
var _ validator.List = allowedValuesListValidator{}

// allowedValuesListValidator validates that each string item in a list matches one of the allowed values.
type allowedValuesListValidator struct {
	allowedValues []string
}

// Description provides a human-readable description of the validator's purpose.
func (v allowedValuesListValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v allowedValuesListValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("each value in the list must be one of: %q", v.allowedValues)
}

// ValidateList performs the validation logic on each item in the list attribute.
func (v allowedValuesListValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	for i, elem := range request.ConfigValue.Elements() {
		strElem, ok := elem.(types.String)
		if !ok || strElem.IsNull() || strElem.IsUnknown() {
			continue
		}

		itemIsValid := false
		for _, allowedValue := range v.allowedValues {
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

// StringListAllowedValues returns a validator that ensures each string in a list matches one of the allowed values.
func StringListAllowedValues(values ...string) validator.List {
	return allowedValuesListValidator{
		allowedValues: values,
	}
}
