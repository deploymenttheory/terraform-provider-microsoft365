package graphBetaWindowsUpdatesAutopatchContentApproval

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchContentApprovalResourceModel struct {
	ID                 types.String        `tfsdk:"id"`
	UpdatePolicyId     types.String        `tfsdk:"update_policy_id"`
	CatalogEntryId     types.String        `tfsdk:"catalog_entry_id"`
	CatalogEntryType   types.String        `tfsdk:"catalog_entry_type"`
	IsRevoked          types.Bool          `tfsdk:"is_revoked"`
	CreatedDateTime    types.String        `tfsdk:"created_date_time"`
	RevokedDateTime    types.String        `tfsdk:"revoked_date_time"`
	DeploymentSettings *DeploymentSettings `tfsdk:"deployment_settings"`
	Timeouts           timeouts.Value      `tfsdk:"timeouts"`
}

type DeploymentSettings struct {
	Schedule *Schedule `tfsdk:"schedule"`
}

type Schedule struct {
	StartDateTime  types.String    `tfsdk:"start_date_time"`
	GradualRollout *GradualRollout `tfsdk:"gradual_rollout"`
}

type GradualRollout struct {
	EndDateTime types.String `tfsdk:"end_date_time"`
}
