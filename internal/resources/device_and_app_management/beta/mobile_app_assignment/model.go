package graphBetaMobileAppAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppAssignmentResource represents the mobile app assignment structure
type MobileAppAssignmentResourceModel struct {
	ID       types.String                     `tfsdk:"id"`
	Intent   types.String                     `tfsdk:"intent"`
	Target   AllLicensedUsersAssignmentTarget `tfsdk:"target"`
	Settings WinGetAppAssignmentSettings      `tfsdk:"settings"`
	Source   types.String                     `tfsdk:"source"`
	SourceID types.String                     `tfsdk:"source_id"`
}

// AllLicensedUsersAssignmentTarget represents the target structure for assignments
type AllLicensedUsersAssignmentTarget struct {
	DeviceAndAppManagementAssignmentFilterID   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
}

// WinGetAppAssignmentSettings represents the settings for the WinGet app assignment
type WinGetAppAssignmentSettings struct {
	Notifications       types.String                 `tfsdk:"notifications"`
	RestartSettings     WinGetAppRestartSettings     `tfsdk:"restart_settings"`
	InstallTimeSettings WinGetAppInstallTimeSettings `tfsdk:"install_time_settings"`
}

// WinGetAppRestartSettings represents the restart settings structure
type WinGetAppRestartSettings struct {
	GracePeriodInMinutes                       types.Int64 `tfsdk:"grace_period_in_minutes"`
	CountdownDisplayBeforeRestartInMinutes     types.Int64 `tfsdk:"countdown_display_before_restart_in_minutes"`
	RestartNotificationSnoozeDurationInMinutes types.Int64 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}

// WinGetAppInstallTimeSettings represents the install time settings structure
type WinGetAppInstallTimeSettings struct {
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
}
