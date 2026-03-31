package graphBetaWindowsUpdatesAutopatchRing

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AssignedGroupModel struct {
	GroupId types.String `tfsdk:"group_id"`
}

type GroupAssignmentModel struct {
	Assignments types.Set `tfsdk:"assignments"`
}

type WindowsUpdatesAutopatchRingResourceModel struct {
	ID                      types.String   `tfsdk:"id"`
	PolicyId                types.String   `tfsdk:"policy_id"`
	DisplayName             types.String   `tfsdk:"display_name"`
	Description             types.String   `tfsdk:"description"`
	IsPaused                types.Bool     `tfsdk:"is_paused"`
	DeferralInDays          types.Int32    `tfsdk:"deferral_in_days"`
	IncludedGroupAssignment types.Object   `tfsdk:"included_group_assignment"`
	ExcludedGroupAssignment types.Object   `tfsdk:"excluded_group_assignment"`
	IsHotpatchEnabled       types.Bool     `tfsdk:"is_hotpatch_enabled"`
	CreatedDateTime         types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime    types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
}
