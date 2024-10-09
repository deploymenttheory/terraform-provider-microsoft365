// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta
package graphBetaMobileAppAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppAssignmentResource represents the mobile app assignment structure
type MobileAppAssignmentResourceModel struct {
	ID       types.String                                  `tfsdk:"id"`
	Intent   types.String                                  `tfsdk:"intent"`
	Target   AllLicensedUsersAssignmentTargetResourceModel `tfsdk:"target"`
	Settings WinGetAppAssignmentSettingsResourceModel      `tfsdk:"settings"`
	Source   types.String                                  `tfsdk:"source"`
	SourceID types.String                                  `tfsdk:"source_id"`
	Timeouts timeouts.Value                                `tfsdk:"timeouts"`
}

// AllLicensedUsersAssignmentTarget represents the target structure for assignments
type AllLicensedUsersAssignmentTargetResourceModel struct {
	DeviceAndAppManagementAssignmentFilterID   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
}

// WinGetAppAssignmentSettings represents the settings for the WinGet app assignment
type WinGetAppAssignmentSettingsResourceModel struct {
	Notifications       types.String                              `tfsdk:"notifications"`
	RestartSettings     WinGetAppRestartSettingsResourceModel     `tfsdk:"restart_settings"`
	InstallTimeSettings WinGetAppInstallTimeSettingsResourceModel `tfsdk:"install_time_settings"`
}

// WinGetAppRestartSettings represents the restart settings structure
type WinGetAppRestartSettingsResourceModel struct {
	GracePeriodInMinutes                       types.Int64 `tfsdk:"grace_period_in_minutes"`
	CountdownDisplayBeforeRestartInMinutes     types.Int64 `tfsdk:"countdown_display_before_restart_in_minutes"`
	RestartNotificationSnoozeDurationInMinutes types.Int64 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}

// WinGetAppInstallTimeSettings represents the install time settings structure
type WinGetAppInstallTimeSettingsResourceModel struct {
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
}
