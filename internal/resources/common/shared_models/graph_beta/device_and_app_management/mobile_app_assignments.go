// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta
package sharedmodels

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MobileAppAssignmentResourceModel struct {
	ID                   types.String          `tfsdk:"id"`
	MobileAppID          types.String          `tfsdk:"mobile_app_id"`
	MobileAppAssignments []MobileAppAssignment `tfsdk:"mobile_app_assignments"`
}

type MobileAppAssignment struct {
	ID       types.String                 `tfsdk:"id"`
	Target   Target                       `tfsdk:"target"`
	Intent   types.String                 `tfsdk:"intent"`
	Settings *WinGetAppAssignmentSettings `tfsdk:"settings"`
	Source   types.String                 `tfsdk:"source"`
	SourceId types.String                 `tfsdk:"source_id"`
}

type Target struct {
	TargetType                                 types.String `tfsdk:"target_type"` // microsoft.graph.groupAssignmentTarget, microsoft.graph.allLicensedUsersAssignmentTarget, etc.
	GroupID                                    types.String `tfsdk:"group_id"`
	DeviceAndAppManagementAssignmentFilterID   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	IsExclusionGroup                           types.Bool   `tfsdk:"is_exclusion_group"`
}

type WinGetAppAssignmentSettings struct {
	Notifications       types.String                  `tfsdk:"notifications"`
	InstallTimeSettings *WinGetAppInstallTimeSettings `tfsdk:"install_time_settings"`
	RestartSettings     *WinGetAppRestartSettings     `tfsdk:"restart_settings"`
}

type WinGetAppInstallTimeSettings struct {
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
}

type WinGetAppRestartSettings struct {
	GracePeriodInMinutes                       types.Int64 `tfsdk:"grace_period_in_minutes"`
	CountdownDisplayBeforeRestartInMinutes     types.Int64 `tfsdk:"countdown_display_before_restart_in_minutes"`
	RestartNotificationSnoozeDurationInMinutes types.Int64 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}

// Constants for various field values
const (
	// Target Types
	TargetTypeGroup               = "microsoft.graph.groupAssignmentTarget"
	TargetTypeExclusionGroup      = "microsoft.graph.exclusionGroupAssignmentTarget"
	TargetTypeAllLicensedUsers    = "microsoft.graph.allLicensedUsersAssignmentTarget"
	TargetTypeAllDevices          = "microsoft.graph.allDevicesAssignmentTarget"
	TargetTypeConfigMgrCollection = "microsoft.graph.configurationManagerCollectionAssignmentTarget"

	// Assignment Filter Types
	FilterTypeNone    = "none"
	FilterTypeInclude = "include"
	FilterTypeExclude = "exclude"

	// Install Intents
	IntentAvailable = "available"
	IntentRequired  = "required"
	IntentUninstall = "uninstall"

	// WinGet App Notifications
	NotificationShowAll    = "showAll"
	NotificationShowReboot = "showReboot"

	// Assignment Sources
	SourceDirect     = "direct"
	SourcePolicySets = "policySets"
)
