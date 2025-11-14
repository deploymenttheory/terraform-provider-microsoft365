package attribute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// jsonSchemaValidator validates that a string is valid JSON and doesn't exceed maximum nesting depth
type jsonSchemaValidator struct {
	maxDepth int
}

// JSONSchemaValidator returns a validator which ensures that a string attribute
// is a valid JSON and doesn't exceed the maximum allowed nesting depth
func JSONSchemaValidator() validator.String {
	return &jsonSchemaValidator{
		maxDepth: 20, // Maximum allowed nesting depth
	}
}

// Description describes the validation in plain text formatting.
func (v jsonSchemaValidator) Description(_ context.Context) string {
	return fmt.Sprintf("JSON string must be valid and not exceed %d levels of nesting", v.maxDepth)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v jsonSchemaValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v jsonSchemaValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// First check if it's valid JSON
	var jsonData any
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &jsonData); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON String",
			fmt.Sprintf("String must be valid JSON, got error: %s", err),
		)
		return
	}

	// Check nesting depth starting from 0
	depth := getJSONDepth(jsonData)
	if depth > v.maxDepth {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"JSON Nesting Too Deep",
			fmt.Sprintf("JSON nesting depth cannot exceed %d levels, got %d levels", v.maxDepth, depth),
		)
	}
}

// getJSONDepth recursively determines the maximum nesting depth of a JSON structure
func getJSONDepth(v any) int {
	switch val := v.(type) {
	case map[string]any:
		maxChildDepth := 0
		for _, child := range val {
			childDepth := getJSONDepth(child)
			if childDepth > maxChildDepth {
				maxChildDepth = childDepth
			}
		}
		return maxChildDepth + 1
	case []any:
		maxChildDepth := 0
		for _, child := range val {
			childDepth := getJSONDepth(child)
			if childDepth > maxChildDepth {
				maxChildDepth = childDepth
			}
		}
		return maxChildDepth + 1
	default:
		return 0
	}
}
