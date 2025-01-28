package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
