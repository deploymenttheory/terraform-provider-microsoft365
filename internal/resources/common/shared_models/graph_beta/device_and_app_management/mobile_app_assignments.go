// Base resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta
package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type MobileAppAssignmentResourceModel struct {
	Id       types.String                             `tfsdk:"id"`
	Intent   types.String                             `tfsdk:"intent"`
	Target   AssignmentTargetResourceModel            `tfsdk:"target"`
	Settings MobileAppAssignmentSettingsResourceModel `tfsdk:"settings"`
	Source   types.String                             `tfsdk:"source"`
	SourceId types.String                             `tfsdk:"source_id"`
}

type AssignmentTargetResourceModel struct {
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"` // allDevicesAssignmentTarget, allLicensedUsersAssignmentTarget, androidFotaDeploymentAssignmentTarget, configurationManagerCollectionAssignmentTarget, exclusionGroupAssignmentTarget, groupAssignmentTarget
}

type MobileAppAssignmentSettingsResourceModel struct {
	AndroidManagedStore       *AndroidManagedStoreAssignmentSettingsResourceModel          `tfsdk:"android_managed_store"`
	IosLob                    *IosLobAppAssignmentSettingsResourceModel                    `tfsdk:"ios_lob"`
	IosStore                  *IosStoreAppAssignmentSettingsResourceModel                  `tfsdk:"ios_store"`
	IosVpp                    *IosVppAppAssignmentSettingsResourceModel                    `tfsdk:"ios_vpp"`
	MacOsLob                  *MacOsLobAppAssignmentSettingsResourceModel                  `tfsdk:"mac_os_lob"`
	MacOsVpp                  *MacOsVppAppAssignmentSettingsResourceModel                  `tfsdk:"mac_os_vpp"`
	MicrosoftStoreForBusiness *MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel `tfsdk:"microsoft_store_for_business"`
	Win32Catalog              *Win32CatalogAppAssignmentSettingsResourceModel              `tfsdk:"win32_catalog"`
	Win32Lob                  *Win32LobAppAssignmentSettingsResourceModel                  `tfsdk:"win32_lob"`
	WindowsAppX               *WindowsAppXAssignmentSettingsResourceModel                  `tfsdk:"windows_app_x"`
	WindowsUniversalAppX      *WindowsUniversalAppXAssignmentSettingsResourceModel         `tfsdk:"windows_universal_app_x"`
	WinGet                    *WinGetAppAssignmentSettingsResourceModel                    `tfsdk:"win_get"`
}

type AndroidManagedStoreAssignmentSettingsResourceModel struct {
	AndroidManagedStoreAppTrackIds types.List   `tfsdk:"android_managed_store_app_track_ids"`
	AutoUpdateMode                 types.String `tfsdk:"auto_update_mode"`
}

type IosLobAppAssignmentSettingsResourceModel struct {
	IsRemovable              types.Bool   `tfsdk:"is_removable"`
	PreventManagedAppBackup  types.Bool   `tfsdk:"prevent_managed_app_backup"`
	UninstallOnDeviceRemoval types.Bool   `tfsdk:"uninstall_on_device_removal"`
	VpnConfigurationId       types.String `tfsdk:"vpn_configuration_id"`
}

type IosStoreAppAssignmentSettingsResourceModel struct {
	IsRemovable              types.Bool   `tfsdk:"is_removable"`
	PreventManagedAppBackup  types.Bool   `tfsdk:"prevent_managed_app_backup"`
	UninstallOnDeviceRemoval types.Bool   `tfsdk:"uninstall_on_device_removal"`
	VpnConfigurationId       types.String `tfsdk:"vpn_configuration_id"`
}

type IosVppAppAssignmentSettingsResourceModel struct {
	IsRemovable              types.Bool   `tfsdk:"is_removable"`
	PreventAutoAppUpdate     types.Bool   `tfsdk:"prevent_auto_app_update"`
	PreventManagedAppBackup  types.Bool   `tfsdk:"prevent_managed_app_backup"`
	UninstallOnDeviceRemoval types.Bool   `tfsdk:"uninstall_on_device_removal"`
	UseDeviceLicensing       types.Bool   `tfsdk:"use_device_licensing"`
	VpnConfigurationId       types.String `tfsdk:"vpn_configuration_id"`
}

type MacOsLobAppAssignmentSettingsResourceModel struct {
	UninstallOnDeviceRemoval types.Bool `tfsdk:"uninstall_on_device_removal"`
}

type MacOsVppAppAssignmentSettingsResourceModel struct {
	PreventAutoAppUpdate     types.Bool `tfsdk:"prevent_auto_app_update"`
	PreventManagedAppBackup  types.Bool `tfsdk:"prevent_managed_app_backup"`
	UninstallOnDeviceRemoval types.Bool `tfsdk:"uninstall_on_device_removal"`
	UseDeviceLicensing       types.Bool `tfsdk:"use_device_licensing"`
}

type MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel struct {
	UseDeviceContext types.Bool `tfsdk:"use_device_context"`
}

type Win32CatalogAppAssignmentSettingsResourceModel struct {
	AutoUpdateSettings           *Win32LobAppAutoUpdateSettingsResourceModel      `tfsdk:"auto_update_settings"`
	DeliveryOptimizationPriority types.String                                     `tfsdk:"delivery_optimization_priority"`
	InstallTimeSettings          *MobileAppInstallTimeSettingsResourceModel       `tfsdk:"install_time_settings"`
	Notifications                types.String                                     `tfsdk:"notifications"`
	RestartSettings              *MobileAppAssignmentSettingsRestartResourceModel `tfsdk:"restart_settings"`
}

type MobileAppInstallTimeSettingsResourceModel struct {
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
	StartDateTime    types.String `tfsdk:"start_date_time"`
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
}

type Win32LobAppAutoUpdateSettingsResourceModel struct {
	AutoUpdateSupersededAppsState types.String `tfsdk:"auto_update_superseded_apps_state"`
}

type MobileAppAssignmentSettingsRestartResourceModel struct {
	GracePeriod                       types.Int32 `tfsdk:"grace_period_in_minutes"`
	CountdownDisplayBeforeRestart     types.Int32 `tfsdk:"countdown_display_before_restart_in_minutes"`
	RestartNotificationSnoozeDuration types.Int32 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}

type MobileAppAssignmentSettingsInstallResourceModel struct {
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
}

type Win32LobAppAssignmentSettingsResourceModel struct {
	AutoUpdateSettings           *Win32LobAppAutoUpdateSettingsResourceModel      `tfsdk:"auto_update_settings"`
	DeliveryOptimizationPriority types.String                                     `tfsdk:"delivery_optimization_priority"`
	InstallTimeSettings          *MobileAppInstallTimeSettingsResourceModel       `tfsdk:"install_time_settings"`
	Notifications                types.String                                     `tfsdk:"notifications"`
	RestartSettings              *MobileAppAssignmentSettingsRestartResourceModel `tfsdk:"restart_settings"`
}

type WindowsAppXAssignmentSettingsResourceModel struct {
	UseDeviceContext types.Bool `tfsdk:"use_device_context"`
}

type WindowsUniversalAppXAssignmentSettingsResourceModel struct {
	UseDeviceContext types.Bool `tfsdk:"use_device_context"`
}

type WinGetAppAssignmentSettingsResourceModel struct {
	InstallTimeSettings *WinGetAppInstallTimeSettingsResourceModel `tfsdk:"install_time_settings"`
	Notifications       types.String                               `tfsdk:"notifications"` // Values: showAll, showReboot, hideAll, unknownFutureValue
	RestartSettings     *WinGetAppRestartSettingsResourceModel     `tfsdk:"restart_settings"`
}

type WinGetAppInstallTimeSettingsResourceModel struct {
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
}

type WinGetAppRestartSettingsResourceModel struct {
	CountdownDisplayBeforeRestartInMinutes     types.Int32 `tfsdk:"countdown_display_before_restart_in_minutes"`
	GracePeriodInMinutes                       types.Int32 `tfsdk:"grace_period_in_minutes"`
	RestartNotificationSnoozeDurationInMinutes types.Int32 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}
