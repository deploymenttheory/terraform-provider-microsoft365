package graphBetaCloudPcProvisioningPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ModifyPlan preserves computed values when the configuration has no effective changes.
func (r *CloudPcProvisioningPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var planned, previous, configured map[string]tftypes.Value
	if req.Plan.Raw.As(&planned) != nil || req.State.Raw.As(&previous) != nil || req.Config.Raw.As(&configured) != nil {
		return
	}
	candidateValues := make(map[string]tftypes.Value, len(planned))
	for name, value := range planned {
		candidateValues[name] = value
	}
	for _, name := range []string{"alternate_resource_url", "cloud_pc_group_display_name", "description", "grace_period_in_hours", "image_display_name"} {
		value, exists := planned[name]
		if !exists || value.IsKnown() {
			continue
		}
		if config, exists := configured[name]; exists && !config.IsNull() {
			continue
		}
		if state, exists := previous[name]; exists {
			candidateValues[name] = state
		}
	}
	candidate := tftypes.NewValue(req.Plan.Raw.Type(), candidateValues)
	// A real configuration change keeps computed outputs unknown until the update completes.
	if candidate.Equal(req.State.Raw) {
		resp.Plan.Raw = candidate
	}
}
