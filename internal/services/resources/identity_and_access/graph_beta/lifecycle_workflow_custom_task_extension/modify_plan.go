package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan sets reply_mode from callback_configuration, which is how Graph derives it: callback
// when a callback configuration is set and none otherwise. This keeps the computed value known at
// plan time, so adding or removing the callback shows the new reply mode and an unchanged
// configuration produces no diff.
func (r *LifecycleWorkflowCustomTaskExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var callback types.Object
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("callback_configuration"), &callback)...)
	if resp.Diagnostics.HasError() || callback.IsUnknown() {
		return
	}

	replyMode := "none"
	if !callback.IsNull() {
		replyMode = "callback"
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("reply_mode"), types.StringValue(replyMode))...)

	tflog.Debug(ctx, "Modify Plan - set reply_mode from callback_configuration")
}
