// REF: https://services.autopatch.microsoft.com/device/v2/autopatchGroups
package graphBetaAutopatchGroups

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AutopatchGroupsResourceModel represents the Terraform resource model
type AutopatchGroupsResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	Name                       types.String   `tfsdk:"name"`
	Description                types.String   `tfsdk:"description"`
	TenantId                   types.String   `tfsdk:"tenant_id"`
	Type                       types.String   `tfsdk:"type"`
	Status                     types.String   `tfsdk:"status"`
	IsLockedByPolicy           types.Bool     `tfsdk:"is_locked_by_policy"`
	DistributionType           types.String   `tfsdk:"distribution_type"`
	ReadOnly                   types.Bool     `tfsdk:"read_only"`
	NumberOfRegisteredDevices  types.Int32    `tfsdk:"number_of_registered_devices"`
	UserHasAllScopeTag         types.Bool     `tfsdk:"user_has_all_scope_tag"`
	FlowId                     types.String   `tfsdk:"flow_id"`
	FlowType                   types.String   `tfsdk:"flow_type"`
	FlowStatus                 types.String   `tfsdk:"flow_status"`
	UmbrellaGroupId            types.String   `tfsdk:"umbrella_group_id"`
	EnableDriverUpdate         types.Bool     `tfsdk:"enable_driver_update"`
	EnabledContentTypes        types.Int32    `tfsdk:"enabled_content_types"`
	GlobalUserManagedAadGroups types.Set      `tfsdk:"global_user_managed_aad_groups"`
	DeploymentGroups           types.List     `tfsdk:"deployment_groups"`
	ScopeTags                  types.Set      `tfsdk:"scope_tags"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`
}

// Terraform resource model nested types
type GlobalUserManagedAadGroup struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type DeploymentGroup struct {
	AadId                         types.String                   `tfsdk:"aad_id"`
	Name                          types.String                   `tfsdk:"name"`
	Distribution                  types.Int32                    `tfsdk:"distribution"`
	FailedPrerequisiteCheckCount  types.Int32                    `tfsdk:"failed_prerequisite_check_count"`
	UserManagedAadGroups          types.Set                      `tfsdk:"user_managed_aad_groups"`
	DeploymentGroupPolicySettings *DeploymentGroupPolicySettings `tfsdk:"deployment_group_policy_settings"`
}

type UserManagedAadGroup struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"` // "Device" or "None" .
}

type DeploymentGroupPolicySettings struct {
	AadGroupName                    types.String                     `tfsdk:"aad_group_name"`
	IsUpdateSettingsModified        types.Bool                       `tfsdk:"is_update_settings_modified"`
	DeviceConfigurationSetting      *DeviceConfigurationSetting      `tfsdk:"device_configuration_setting"`
	DnfUpdateCloudSetting           *DnfUpdateCloudSetting           `tfsdk:"dnf_update_cloud_setting"`
	OfficeDCv2Setting               *OfficeDCv2Setting               `tfsdk:"office_dcv2_setting"`
	EdgeDCv2Setting                 *EdgeDCv2Setting                 `tfsdk:"edge_dcv2_setting"`
	FeatureUpdateAnchorCloudSetting *FeatureUpdateAnchorCloudSetting `tfsdk:"feature_update_anchor_cloud_setting"`
}

type DeviceConfigurationSetting struct {
	PolicyId                  types.String               `tfsdk:"policy_id"`
	UpdateBehavior            types.String               `tfsdk:"update_behavior"`
	NotificationSetting       types.String               `tfsdk:"notification_setting"`
	QualityDeploymentSettings *QualityDeploymentSettings `tfsdk:"quality_deployment_settings"`
	FeatureDeploymentSettings *FeatureDeploymentSettings `tfsdk:"feature_deployment_settings"`
}

type QualityDeploymentSettings struct {
	Deadline    types.Int32 `tfsdk:"deadline"`
	Deferral    types.Int32 `tfsdk:"deferral"`
	GracePeriod types.Int32 `tfsdk:"grace_period"`
}

type FeatureDeploymentSettings struct {
	Deadline types.Int32 `tfsdk:"deadline"`
	Deferral types.Int32 `tfsdk:"deferral"`
}

type DnfUpdateCloudSetting struct {
	PolicyId                 types.String `tfsdk:"policy_id"`
	ApprovalType             types.String `tfsdk:"approval_type"`
	DeploymentDeferralInDays types.Int32  `tfsdk:"deployment_deferral_in_days"`
}

type OfficeDCv2Setting struct {
	PolicyId                types.String `tfsdk:"policy_id"`
	Deadline                types.Int32  `tfsdk:"deadline"`
	Deferral                types.Int32  `tfsdk:"deferral"`
	HideUpdateNotifications types.Bool   `tfsdk:"hide_update_notifications"`
	TargetChannel           types.String `tfsdk:"target_channel"`
	EnableAutomaticUpdate   types.Bool   `tfsdk:"enable_automatic_update"`
	HideEnableDisableUpdate types.Bool   `tfsdk:"hide_enable_disable_update"`
	EnableOfficeMgmt        types.Bool   `tfsdk:"enable_office_mgmt"`
	UpdatePath              types.String `tfsdk:"update_path"`
}

type EdgeDCv2Setting struct {
	PolicyId      types.String `tfsdk:"policy_id"`
	TargetChannel types.String `tfsdk:"target_channel"`
}

type FeatureUpdateAnchorCloudSetting struct {
	TargetOSVersion                                   types.String `tfsdk:"target_os_version"`
	InstallLatestWindows10OnWindows11IneligibleDevice types.Bool   `tfsdk:"install_latest_windows10_on_windows11_ineligible_device"`
	PolicyId                                          types.String `tfsdk:"policy_id"`
}
