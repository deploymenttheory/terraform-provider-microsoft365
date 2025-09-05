package graphBetaAppControlForBusinessPolicy

import (
	"context"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan handles plan modification for diff suppression
func (r *AppControlForBusinessPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	tflog.Debug(ctx, "Modifying plan for app control for business policy")

	// Extract the models from plan and state
	var plan, state AppControlForBusinessPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle XML content normalization to suppress BOM differences
	if !plan.PolicyXML.IsNull() && !state.PolicyXML.IsNull() {
		planXML := plan.PolicyXML.ValueString()
		stateXML := state.PolicyXML.ValueString()

		// Check if the XML likely came from a file
		fromFile := normalize.LikelyFromFile(planXML)

		// Normalize both by removing BOM if present
		planXMLNormalized := strings.TrimPrefix(planXML, "\ufeff")
		stateXMLNormalized := strings.TrimPrefix(stateXML, "\ufeff")

		// If they're the same after normalization
		if planXMLNormalized == stateXMLNormalized {
			// For file-sourced XML, keep the BOM if it was present
			if fromFile {
				// Use state value which should have BOM preserved by ReverseNormalizeXMLContent
				plan.PolicyXML = state.PolicyXML
				tflog.Debug(ctx, "Suppressing diff for file-sourced policy_xml due to BOM differences only")
			} else {
				// For inline XML, ensure we don't have a BOM
				if strings.HasPrefix(stateXML, "\ufeff") {
					// Use plan value which shouldn't have BOM
					tflog.Debug(ctx, "Suppressing diff for inline policy_xml due to BOM differences only")
				} else {
					// Both don't have BOM, use state to prevent unnecessary updates
					plan.PolicyXML = state.PolicyXML
				}
			}
		}
	}

	// Set the modified plan
	resp.Plan.Set(ctx, plan)
}
