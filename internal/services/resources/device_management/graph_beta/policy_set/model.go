// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-policyset-policyset?view=graph-rest-beta
package graphBetaPolicySet

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PolicySetResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	Status               types.String   `tfsdk:"status"`
	ErrorCode            types.String   `tfsdk:"error_code"`
	RoleScopeTagIds      types.Set      `tfsdk:"role_scope_tag_ids"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Assignments          types.Set      `tfsdk:"assignments"`
	Items                types.Set      `tfsdk:"items"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// Use the common assignment model
type PolicySetAssignmentModel = sharedmodels.InclusionGroupAndExclusionGroupAssignmentModel

type PolicySetItemModel struct {
	ID        types.String `tfsdk:"id"`
	PayloadId types.String `tfsdk:"payload_id"`
	Type      types.String `tfsdk:"type"`
	Intent    types.String `tfsdk:"intent"`
	Settings  types.Object `tfsdk:"settings"`
}

type PolicySetItemSettingsModel struct {
	ODataType                types.String `tfsdk:"odata_type"`
	VpnConfigurationId       types.String `tfsdk:"vpn_configuration_id"`
	UninstallOnDeviceRemoval types.Bool   `tfsdk:"uninstall_on_device_removal"`
	IsRemovable              types.Bool   `tfsdk:"is_removable"`
	PreventManagedAppBackup  types.Bool   `tfsdk:"prevent_managed_app_backup"`
}
