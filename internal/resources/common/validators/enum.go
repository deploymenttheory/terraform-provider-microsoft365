package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure that the struct implements the required interface
var _ validator.String = enumValidator{}

// enumValidator validates that the value matches one of the expected enum values.
type enumValidator struct {
	allowedEnumValues []string
}

// Description provides a human-readable description of the validator's purpose.
func (v enumValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v enumValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("all individual values must be one of: %q", v.allowedEnumValues)
}

// ValidateString performs the validation logic on a string attribute.
func (v enumValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	configValue := request.ConfigValue
	if configValue.IsNull() || configValue.IsUnknown() || configValue.ValueString() == "" {
		return
	}

	// Split and validate each comma-separated value
	for _, configItem := range strings.Split(configValue.ValueString(), ",") {
		itemIsValid := false
		for _, allowedValue := range v.allowedEnumValues {
			if strings.TrimSpace(configItem) == allowedValue {
				itemIsValid = true
				break
			}
		}
		if !itemIsValid {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				request.Path,
				v.Description(ctx),
				configValue.String(),
			))
		}
	}
}

// EnumValues returns an instance of the validator with specified allowed values.
func EnumValues(values ...string) validator.String {
	return enumValidator{
		allowedEnumValues: values,
	}
}
