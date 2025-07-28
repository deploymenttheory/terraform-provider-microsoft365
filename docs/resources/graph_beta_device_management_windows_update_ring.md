---
page_title: "microsoft365_graph_beta_device_management_windows_update_ring Resource - terraform-provider-microsoft365"
subcategory: "Device Management"
description: |-
  Manages Windows Update for Business configuration policies using the /deviceManagement/deviceConfigurations endpoint. This resource controls Windows Update settings including feature update deferrals, quality update schedules, driver management, and restart behaviors for managed Windows 10/11 devices.
---

# microsoft365_graph_beta_device_management_windows_update_ring (Resource)

Manages Windows Update for Business configuration policies using the `/deviceManagement/deviceConfigurations` endpoint. This resource controls Windows Update settings including feature update deferrals, quality update schedules, driver management, and restart behaviors for managed Windows 10/11 devices.

## Microsoft Documentation

- [Windows Update for Business Configuration](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdateforbusinessconfiguration?view=graph-rest-beta)
- [Windows Update Settings](https://learn.microsoft.com/en-us/mem/intune/protect/windows-update-settings)
- [Windows Update for Business Overview](https://learn.microsoft.com/en-us/windows/deployment/update/waas-manage-updates-wufb)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Read**: `DeviceManagementConfiguration.Read.All`
- **Write**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.23.0  | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Basic Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "basic_update_ring" {
  display_name       = "Standard Windows Update Ring"
  description        = "Default update ring for standard workstations"
  role_scope_tag_ids = ["0"]

  microsoft_update_service_allowed             = true
  drivers_excluded                             = false
  quality_updates_deferral_period_in_days      = 7
  feature_updates_deferral_period_in_days      = 14
  allow_windows11_upgrade                      = false
  quality_updates_paused                       = false
  feature_updates_paused                       = false
  business_ready_updates_only                  = "businessReadyOnly"
  automatic_update_mode                        = "autoInstallAtMaintenanceTime"
  active_hours_start                           = "08:00:00"
  active_hours_end                             = "17:00:00"
  user_pause_access                            = "enabled"
  user_windows_update_scan_access              = "enabled"
  update_notification_level                    = "defaultNotifications"
  deadline_for_feature_updates_in_days         = 7
  deadline_for_quality_updates_in_days         = 3
  deadline_grace_period_in_days                = 2
  skip_checks_before_restart                   = false
  postpone_reboot_until_after_deadline         = true
  engaged_restart_deadline_in_days             = 7
  engaged_restart_snooze_schedule_in_days      = 2
  engaged_restart_transition_schedule_in_days  = 7
  auto_restart_notification_dismissal          = "notConfigured"
  schedule_restart_warning_in_hours            = 4
  schedule_imminent_restart_warning_in_minutes = 30
  delivery_optimization_mode                   = "httpWithPeeringNat"
  prerelease_features                          = "notAllowed"
  update_weeks                                 = "everyWeek"
  feature_updates_rollback_window_in_days      = 10

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}

# Example 2: Minimal Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "minimal_update_ring" {
  display_name = "Minimal Windows Update Ring"
  description  = "Basic update ring with minimal configuration"

  # Only required fields
  automatic_update_mode = "autoInstallAndRebootAtScheduledTime"
}

# Example 3: Advanced Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "advanced_update_ring" {
  display_name = "Advanced Windows Update Ring"
  description  = "Advanced update ring for specialized workstations"

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Update service configuration
  microsoft_update_service_allowed = true
  drivers_excluded                 = true

  # Update deferral configuration
  quality_updates_deferral_period_in_days = 14
  feature_updates_deferral_period_in_days = 30

  # Windows 11 upgrade settings
  allow_windows11_upgrade = true

  # Update pause configuration
  quality_updates_paused = false
  feature_updates_paused = false

  # Update branch configuration
  business_ready_updates_only = "businessReadyOnly"

  # Automatic update mode
  automatic_update_mode = "autoInstallAndRebootAtScheduledTime"

  # Active hours configuration
  active_hours_start = "07:00:00"
  active_hours_end   = "19:00:00"

  # User control settings
  user_pause_access               = "disabled"
  user_windows_update_scan_access = "enabled"
  update_notification_level       = "restartWarningsOnly"

  # Deadline configuration
  deadline_for_feature_updates_in_days = 14
  deadline_for_quality_updates_in_days = 7
  deadline_grace_period_in_days        = 3

  # Restart behavior
  skip_checks_before_restart           = false
  postpone_reboot_until_after_deadline = true

  # Engaged restart settings
  engaged_restart_deadline_in_days            = 14
  engaged_restart_snooze_schedule_in_days     = 3
  engaged_restart_transition_schedule_in_days = 14

  # Restart notifications
  auto_restart_notification_dismissal          = "automatic"
  schedule_restart_warning_in_hours            = 8
  schedule_imminent_restart_warning_in_minutes = 60

  # Delivery optimization
  delivery_optimization_mode = "httpWithInternetPeering"

  # Feature management
  prerelease_features = "notAllowed"
  update_weeks        = "firstWeek"

  # Rollback settings
  feature_updates_rollback_window_in_days = 20
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Admin provided name of the device configuration. Inherited from deviceConfiguration.

### Optional

- `active_hours_end` (String) Active Hours End. Part of the Installation Schedule.
- `active_hours_start` (String) Active Hours Start. Part of the Installation Schedule.
- `additional_properties` (Map of String) Additional properties that are not yet exposed in the API.
- `allow_windows11_upgrade` (Boolean) When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.
- `auto_restart_notification_dismissal` (String) Specify the method by which the auto-restart required notification is dismissed. Possible values are: NotConfigured, Automatic, User. Returned by default. Query parameters are not supported. Possible values are: notConfigured, automatic, user, unknownFutureValue.
- `automatic_update_mode` (String) The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime, AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, WindowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. Possible values are: userDefined, notifyDownload, autoInstallAtMaintenanceTime, autoInstallAndRebootAtMaintenanceTime, autoInstallAndRebootAtScheduledTime, autoInstallAndRebootWithoutEndUserControl.
- `business_ready_updates_only` (String) Determines which branch devices will receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, WindowsInsiderBuildSlow, WindowsInsiderBuildRelease. Returned by default. Query parameters are not supported. Possible values are: userDefined, all, businessReadyOnly, windowsInsiderBuildFast, windowsInsiderBuildSlow, windowsInsiderBuildRelease.
- `deadline_for_feature_updates_in_days` (Number) Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_for_quality_updates_in_days` (Number) Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_grace_period_in_days` (Number) Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.
- `delivery_optimization_mode` (String) The Delivery Optimization Mode. Possible values are: UserDefined, HttpOnly, HttpWithPeeringNat, HttpWithPeeringPrivateGroup, HttpWithInternetPeering, SimpleDownload, BypassMode. UserDefined allows the user to set. Returned by default. Query parameters are not supported. Possible values are: userDefined, httpOnly, httpWithPeeringNat, httpWithPeeringPrivateGroup, httpWithInternetPeering, simpleDownload, bypassMode.
- `description` (String) Admin provided description of the Device Configuration. Inherited from deviceConfiguration.
- `drivers_excluded` (Boolean) When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.
- `engaged_restart_deadline_in_days` (Number) Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.
- `engaged_restart_snooze_schedule_for_feature_updates_in_days` (Number) Number of days a user can snooze Engaged Restart reminder notifications for feature updates.
- `engaged_restart_snooze_schedule_in_days` (Number) Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.
- `engaged_restart_transition_schedule_for_feature_updates_in_days` (Number) Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart for feature updates.
- `engaged_restart_transition_schedule_in_days` (Number) Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `feature_updates_deferral_period_in_days` (Number) Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `feature_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s
- `feature_updates_rollback_window_in_days` (Number) The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.
- `microsoft_update_service_allowed` (Boolean) When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.
- `postpone_reboot_until_after_deadline` (Boolean) When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.
- `prerelease_features` (String) The Pre-Release Features. Possible values are: UserDefined, SettingsOnly, SettingsAndExperimentations, NotAllowed. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. Possible values are: userDefined, settingsOnly, settingsAndExperimentations, notAllowed.
- `quality_updates_deferral_period_in_days` (Number) Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `quality_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `schedule_imminent_restart_warning_in_minutes` (Number) Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.
- `schedule_restart_warning_in_hours` (Number) Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.
- `skip_checks_before_restart` (Boolean) When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `update_notification_level` (String) Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported. Possible values are: notConfigured, defaultNotifications, restartWarningsOnly, disableAllNotifications, unknownFutureValue.
- `update_weeks` (String) Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. Possible values are: userDefined, firstWeek, secondWeek, thirdWeek, fourthWeek, everyWeek, unknownFutureValue.
- `user_pause_access` (String) Specifies whether to enable end user's access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.
- `user_windows_update_scan_access` (String) Specifies whether to disable user's access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.

### Read-Only

- `id` (String) Key of the entity. Inherited from deviceConfiguration.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Windows Update Rings**: This resource manages Windows Update for Business configuration policies, allowing you to control how and when Windows updates are applied to managed devices.
- **Active Hours**: The active hours configuration (start and end time) defines when devices should not automatically restart after updates. Both values must be provided together.
- **Update Deferrals**: Quality updates can be deferred up to 30 days, and feature updates can be deferred up to 30 days.
- **Windows 11 Upgrade**: The `allow_windows11_upgrade` attribute controls whether eligible Windows 10 devices can upgrade to Windows 11.
- **Automatic Update Mode**: Controls how updates are installed and when restarts occur. Options range from notification-only to fully automated installation and restart.
- **Delivery Optimization**: Controls how updates are downloaded and distributed within your network to optimize bandwidth usage.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using group object ID
terraform import microsoft365_graph_beta_device_management_windows_update_ring.example 00000000-0000-0000-0000-000000000000
``` 