// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcusersetting?view=graph-rest-beta
package graphBetaCloudPcUserSetting

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcUserSettingResourceModel struct {
	ID                                 types.String                             `tfsdk:"id"`
	DisplayName                        types.String                             `tfsdk:"display_name"`
	CreatedDateTime                    types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime               types.String                             `tfsdk:"last_modified_date_time"`
	LocalAdminEnabled                  types.Bool                               `tfsdk:"local_admin_enabled"`
	ResetEnabled                       types.Bool                               `tfsdk:"reset_enabled"`
	SelfServiceEnabled                 types.Bool                               `tfsdk:"self_service_enabled"`
	RestorePointSetting                *RestorePointSettingModel                `tfsdk:"restore_point_setting"`
	CrossRegionDisasterRecoverySetting *CrossRegionDisasterRecoverySettingModel `tfsdk:"cross_region_disaster_recovery_setting"`
	NotificationSetting                *NotificationSettingModel                `tfsdk:"notification_setting"`
	Assignments                        []CloudPcUserSettingAssignmentModel      `tfsdk:"assignments"`
	Timeouts                           timeouts.Value                           `tfsdk:"timeouts"`
}

type RestorePointSettingModel struct {
	FrequencyInHours   types.Int32  `tfsdk:"frequency_in_hours"`
	FrequencyType      types.String `tfsdk:"frequency_type"`
	UserRestoreEnabled types.Bool   `tfsdk:"user_restore_enabled"`
}

type CrossRegionDisasterRecoverySettingModel struct {
	MaintainCrossRegionRestorePointEnabled types.Bool                           `tfsdk:"maintain_cross_region_restore_point_enabled"`
	UserInitiatedDisasterRecoveryAllowed   types.Bool                           `tfsdk:"user_initiated_disaster_recovery_allowed"`
	DisasterRecoveryType                   types.String                         `tfsdk:"disaster_recovery_type"`
	DisasterRecoveryNetworkSetting         *DisasterRecoveryNetworkSettingModel `tfsdk:"disaster_recovery_network_setting"`
}

type DisasterRecoveryNetworkSettingModel struct {
	NetworkType types.String `tfsdk:"network_type"`
	RegionName  types.String `tfsdk:"region_name"`
	RegionGroup types.String `tfsdk:"region_group"`
}

type NotificationSettingModel struct {
	RestartPromptsDisabled types.Bool `tfsdk:"restart_prompts_disabled"`
}

// CloudPcUserSettingAssignmentModel represents an assignment of a Cloud PC user setting to a group
type CloudPcUserSettingAssignmentModel struct {
	ID      types.String `tfsdk:"id"`
	GroupId types.String `tfsdk:"group_id"`
}
