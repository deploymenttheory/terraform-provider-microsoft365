// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofile?view=graph-rest-beta
package graphBetaWindowsDriverUpdateProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsDriverUpdateProfileResourceModel struct {
	ID                       types.String   `tfsdk:"id"`
	DisplayName              types.String   `tfsdk:"display_name"`
	Description              types.String   `tfsdk:"description"`
	ApprovalType             types.String   `tfsdk:"approval_type"`
	DeviceReporting          types.Int32    `tfsdk:"device_reporting"`
	NewUpdates               types.Int32    `tfsdk:"new_updates"`
	DeploymentDeferralInDays types.Int32    `tfsdk:"deployment_deferral_in_days"`
	CreatedDateTime          types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime     types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds          types.Set      `tfsdk:"role_scope_tag_ids"`
	InventorySyncStatus      types.Object   `tfsdk:"inventory_sync_status"`
	Assignments              types.Set      `tfsdk:"assignments"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}
