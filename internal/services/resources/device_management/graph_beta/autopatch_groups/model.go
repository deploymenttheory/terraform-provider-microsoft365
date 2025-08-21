// REF: https://services.autopatch.microsoft.com/device/v2/autopatchGroups
package graphBetaAutopatchGroups

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AutopatchGroupsResourceModel represents the Terraform resource model
type AutopatchGroupsResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	Name                         types.String   `tfsdk:"name"`
	Description                  types.String   `tfsdk:"description"`
	TenantId                     types.String   `tfsdk:"tenant_id"`
	Type                         types.String   `tfsdk:"type"`
	Status                       types.String   `tfsdk:"status"`
	IsLockedByPolicy             types.Bool     `tfsdk:"is_locked_by_policy"`
	DistributionType             types.String   `tfsdk:"distribution_type"`
	ReadOnly                     types.Bool     `tfsdk:"read_only"`
	NumberOfRegisteredDevices    types.Int64    `tfsdk:"number_of_registered_devices"`
	UserHasAllScopeTag           types.Bool     `tfsdk:"user_has_all_scope_tag"`
	FlowId                       types.String   `tfsdk:"flow_id"`
	FlowType                     types.String   `tfsdk:"flow_type"`
	FlowStatus                   types.String   `tfsdk:"flow_status"`
	UmbrellaGroupId              types.String   `tfsdk:"umbrella_group_id"`
	EnableDriverUpdate           types.Bool     `tfsdk:"enable_driver_update"`
	EnabledContentTypes          types.Int64    `tfsdk:"enabled_content_types"`
	GlobalUserManagedAadGroups   types.Set      `tfsdk:"global_user_managed_aad_groups"`
	DeploymentGroups             types.Set      `tfsdk:"deployment_groups"`
	ScopeTags                    types.Set      `tfsdk:"scope_tags"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}

// Terraform resource model nested types
type GlobalUserManagedAadGroup struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type DeploymentGroup struct {
	AadId                         types.String                     `tfsdk:"aad_id"`
	Name                          types.String                     `tfsdk:"name"`
	Distribution                  types.Int64                      `tfsdk:"distribution"`
	FailedPrerequisiteCheckCount  types.Int64                      `tfsdk:"failed_prerequisite_check_count"`
	UserManagedAadGroups          types.Set                        `tfsdk:"user_managed_aad_groups"`
	DeploymentGroupPolicySettings *DeploymentGroupPolicySettings   `tfsdk:"deployment_group_policy_settings"`
}

type UserManagedAadGroup struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.Int64  `tfsdk:"type"`
}

type DeploymentGroupPolicySettings struct {
	AadGroupName               types.String                   `tfsdk:"aad_group_name"`
	IsUpdateSettingsModified   types.Bool                     `tfsdk:"is_update_settings_modified"`
	DeviceConfigurationSetting *DeviceConfigurationSetting   `tfsdk:"device_configuration_setting"`
}

type DeviceConfigurationSetting struct {
	PolicyId                   types.String                   `tfsdk:"policy_id"`
	UpdateBehavior             types.String                   `tfsdk:"update_behavior"`
	NotificationSetting        types.String                   `tfsdk:"notification_setting"`
	QualityDeploymentSettings  *QualityDeploymentSettings     `tfsdk:"quality_deployment_settings"`
	FeatureDeploymentSettings  *FeatureDeploymentSettings     `tfsdk:"feature_deployment_settings"`
}

type QualityDeploymentSettings struct {
	Deadline    types.Int64 `tfsdk:"deadline"`
	Deferral    types.Int64 `tfsdk:"deferral"`
	GracePeriod types.Int64 `tfsdk:"grace_period"`
}

type FeatureDeploymentSettings struct {
	Deadline types.Int64 `tfsdk:"deadline"`
	Deferral types.Int64 `tfsdk:"deferral"`
}