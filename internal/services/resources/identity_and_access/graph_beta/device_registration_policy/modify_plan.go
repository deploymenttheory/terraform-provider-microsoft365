package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan reports, at plan time, a change to azure_ad_registration.allowed_to_register when the
// tenant policy in state says the setting is not admin configurable (Microsoft Intune is enabled).
// On create there is no state yet, so the same check runs at apply time against a fresh read.
func (r *DeviceRegistrationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state DeviceRegistrationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isAdminConfigurable, ok := priorAttribute(state.AzureADRegistration, "is_admin_configurable").(basetypes.BoolValue)
	if !ok || isAdminConfigurable.IsNull() || isAdminConfigurable.IsUnknown() || isAdminConfigurable.ValueBool() {
		tflog.Debug(ctx, "Modify Plan - azure_ad_registration is admin configurable, no checks needed")
		return
	}

	planned, ok := priorAttribute(plan.AzureADRegistration, "allowed_to_register").(basetypes.ObjectValue)
	if !ok || planned.IsNull() || planned.IsUnknown() {
		return
	}
	current, ok := priorAttribute(state.AzureADRegistration, "allowed_to_register").(basetypes.ObjectValue)
	if !ok || current.IsNull() || current.IsUnknown() {
		return
	}

	if !membershipsEqual(planned, current) {
		resp.Diagnostics.AddAttributeError(
			path.Root("azure_ad_registration").AtName("allowed_to_register"),
			"azure_ad_registration cannot be changed",
			registrationNotConfigurableDetail,
		)
	}
}
