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
