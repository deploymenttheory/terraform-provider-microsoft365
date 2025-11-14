// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreappconfiguration?view=graph-rest-beta
package graphBetaAndroidManagedDeviceAppConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AndroidManagedDeviceAppConfigurationPolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	Description          types.String   `tfsdk:"description"`
	TargetedMobileApps   types.Set      `tfsdk:"targeted_mobile_apps"`
	RoleScopeTagIds      types.Set      `tfsdk:"role_scope_tag_ids"`
	Version              types.Int32    `tfsdk:"version"`
	PackageId            types.String   `tfsdk:"package_id"`
	PayloadJson          types.String   `tfsdk:"payload_json"`
	ProfileApplicability types.String   `tfsdk:"profile_applicability"`
	ConnectedAppsEnabled types.Bool     `tfsdk:"connected_apps_enabled"`
	AppSupportsOemConfig types.Bool     `tfsdk:"app_supports_oem_config"`
	PermissionActions    types.Set      `tfsdk:"permission_actions"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type AndroidPermissionActionModel struct {
	Permission types.String `tfsdk:"permission"`
	Action     types.String `tfsdk:"action"`
}
