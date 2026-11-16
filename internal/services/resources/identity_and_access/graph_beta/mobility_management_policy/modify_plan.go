package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan plans the computed values the write will produce so applies stay consistent:
//   - applies_to none or all removes every included group, so included_groups is planned as null
//     instead of keeping the prior state value.
//   - is_mdm_enrollment_during_registration_disabled only exists on MDM policies, so it is planned as null
//     for app policies.
func (r *MobilityManagementPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan MobilityManagementPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appliesTo := plan.AppliesTo.ValueString()
	if isConfigured(plan.AppliesTo) && (appliesTo == appliesToNone || appliesTo == appliesToAll) && !plan.IncludedGroups.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("Modify Plan - applies_to is %s, planning included_groups as null", appliesTo))
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("included_groups"), types.SetNull(types.StringType))...)
	}

	if plan.PolicyType.ValueString() == policyTypeApp && !plan.IsMdmEnrollmentDuringRegistrationDisabled.IsNull() {
		tflog.Debug(ctx, "Modify Plan - app policies have no is_mdm_enrollment_during_registration_disabled, planning it as null")
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("is_mdm_enrollment_during_registration_disabled"), types.BoolNull())...)
	}
}
