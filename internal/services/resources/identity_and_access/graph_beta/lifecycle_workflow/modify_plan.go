package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan sets each task's execution_sequence from its position in the tasks list, which is
// how Graph numbers tasks. This keeps the computed value known at plan time, so reordering tasks
// shows the new sequence and an unchanged list produces no diff.
func (r *LifecycleWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var tasksList types.List
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tasks"), &tasksList)...)
	if resp.Diagnostics.HasError() || tasksList.IsNull() || tasksList.IsUnknown() {
		return
	}

	for i, element := range tasksList.Elements() {
		if element.IsUnknown() {
			continue
		}
		sequencePath := path.Root("tasks").AtListIndex(i).AtName("execution_sequence")
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, sequencePath, types.Int64Value(int64(i+1)))...)
	}

	tflog.Debug(ctx, "Modify Plan - set task execution_sequence from list order")
}
