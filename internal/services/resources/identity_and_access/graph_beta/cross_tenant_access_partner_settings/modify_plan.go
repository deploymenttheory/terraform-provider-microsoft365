package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModifyPlan makes inheritance an explicit desired state while retaining the
// existing Optional+Computed behaviour for users who do not opt in.
func (r *CrossTenantAccessPartnerSettingsResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var enforce types.Bool
	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, path.Root("enforce_default_inheritance"), &enforce)...)
	if resp.Diagnostics.HasError() || (!enforce.IsUnknown() && !enforce.ValueBool()) {
		return
	}
	for _, setting := range inheritablePartnerSettings {
		var configured types.Object
		resp.Diagnostics.Append(
			req.Config.GetAttribute(ctx, path.Root(setting.terraformName), &configured)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if configured.IsNull() {
			value := types.ObjectNull(configured.AttributeTypes(ctx))
			if enforce.IsUnknown() {
				// Leave room for either inheritance or the remote value until
				// an expression supplying the enforcement flag is resolved.
				value = types.ObjectUnknown(configured.AttributeTypes(ctx))
			}
			resp.Diagnostics.Append(
				resp.Plan.SetAttribute(ctx, path.Root(setting.terraformName), value)...)
		}
	}
}
