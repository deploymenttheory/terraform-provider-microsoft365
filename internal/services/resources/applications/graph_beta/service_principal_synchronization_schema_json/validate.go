package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = jsonArrayValidator{}

// jsonArrayValidator validates that a string attribute holds a JSON array.
type jsonArrayValidator struct{}

// Description describes the validation in plain text formatting.
func (v jsonArrayValidator) Description(_ context.Context) string {
	return "value must be a JSON array"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v jsonArrayValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v jsonArrayValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	decoded, err := decodeJSON([]byte(req.ConfigValue.ValueString()))
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON", fmt.Sprintf("The value must be a JSON array: %s", err.Error()))
		return
	}

	if _, ok := decoded.([]any); !ok {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON", fmt.Sprintf("The value must be a JSON array, got a JSON %s.", jsonKind(decoded)))
	}
}

// jsonKind names the JSON kind of a decoded value for error messages.
func jsonKind(value any) string {
	switch value.(type) {
	case map[string]any:
		return "object"
	case string:
		return "string"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "number"
	}
}
