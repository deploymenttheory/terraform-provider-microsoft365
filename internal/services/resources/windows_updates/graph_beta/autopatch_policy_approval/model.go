package graphBetaWindowsUpdatesAutopatchPolicyApproval

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdatesAutopatchPolicyApprovalResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	PolicyId           types.String   `tfsdk:"policy_id"`
	CatalogEntryId     types.String   `tfsdk:"catalog_entry_id"`
	Status             types.String   `tfsdk:"status"`
	CreatedDateTime    types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}
