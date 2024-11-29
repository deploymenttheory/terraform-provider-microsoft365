package planmodifiers

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PlanModifyString normalizes the JSON string during the plan phase.
// NormalizeJSONPlanModifier is a custom plan modifier that normalizes JSON strings.
type NormalizeJSONPlanModifier struct{}

// PlanModifyString normalizes the JSON string during the plan phase.
func (m NormalizeJSONPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	jsonInput := req.ConfigValue.ValueString()
	normalized, err := normalize.JSONAlphabetically(jsonInput)
	if err != nil {
		tflog.Warn(ctx, "Failed to normalize JSON in plan modifier", map[string]interface{}{
			"error": err.Error(),
			"value": jsonInput,
		})
		resp.Diagnostics.AddError(
			"JSON Normalization Failed",
			"Failed to normalize the settings JSON string: "+err.Error(),
		)
		return
	}

	resp.PlanValue = types.StringValue(normalized)
	tflog.Debug(ctx, "Normalized JSON during plan phase", map[string]interface{}{
		"original":   jsonInput,
		"normalized": normalized,
	})
}

// Description provides a description of the Plan Modifier.
func (m NormalizeJSONPlanModifier) Description(ctx context.Context) string {
	return "Normalizes JSON strings by sorting object keys alphabetically."
}

// MarkdownDescription provides a markdown-compatible description.
func (m NormalizeJSONPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Normalizes JSON strings by sorting object keys alphabetically."
}
