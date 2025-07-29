package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan modifies the plan for the resource.
func (r *LinuxDeviceCompliancePolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No modifications needed at this time
}
