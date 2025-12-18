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
| v0.39.0  | Preview | Refactored resource to align with api changes and added full scenario based test harness|

## Example Usage

### Scenario 1: Notify Download

This configuration notifies users when updates are available but requires manual download and installation.

```terraform
# Scenario 1: Notify Download
# This configuration notifies users when updates are available for download but does not
# automatically download or install them. Users maintain full control over the update process.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "notify_download" {
  display_name                            = "Windows Update Ring - Notify Download"
  description                             = "Notify users when updates are available for download"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "notifyDownload"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 2: Auto Install at Maintenance Time

This configuration automatically installs updates outside of active hours but requires user interaction to restart.

```terraform
# Scenario 2: Auto Install at Maintenance Time
# This configuration automatically installs updates outside of active hours but requires
# user interaction to restart. Updates install during maintenance windows (outside active hours).

resource "microsoft365_graph_beta_device_management_windows_update_ring" "auto_install_maintenance" {
  display_name                            = "Windows Update Ring - Auto Install at Maintenance Time"
  description                             = "Automatically install updates at maintenance time"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAtMaintenanceTime"
  active_hours_start                      = "08:00:00"
  active_hours_end                        = "17:00:00"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 3: Auto Install and Reboot at Maintenance Time

This configuration automatically installs updates and restarts devices outside of active hours.

```terraform
# Scenario 3: Auto Install and Reboot at Maintenance Time
# This configuration automatically installs updates and restarts devices outside of active hours.
# Devices will automatically reboot during maintenance windows without user interaction.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "auto_reboot_maintenance" {
  display_name                            = "Windows Update Ring - Auto Install and Reboot at Maintenance"
  description                             = "Automatically install and reboot at maintenance time"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAndRebootAtMaintenanceTime"
  active_hours_start                      = "08:00:00"
  active_hours_end                        = "17:00:00"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 4: Auto Install and Restart at Scheduled Time

This configuration automatically installs updates and restarts devices at a specific scheduled time.

```terraform
# Scenario 4: Auto Install and Restart at Scheduled Time
# This configuration automatically installs updates and restarts devices at a specific scheduled
# time and day. Use this for predictable maintenance windows.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "scheduled_install" {
  display_name                            = "Windows Update Ring - Scheduled Install and Restart"
  description                             = "Automatically install and restart at scheduled time"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAndRebootAtScheduledTime"
  scheduled_install_day                   = "everyday"
  scheduled_install_time                  = "03:00:00"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  update_weeks                            = "everyWeek"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 5: Auto Install and Reboot Without End User Control

This configuration provides the most aggressive update policy with no end-user control.

```terraform
# Scenario 5: Auto Install and Reboot Without End User Control
# This configuration provides the most aggressive update policy, automatically installing and
# restarting devices without user interaction or the ability to postpone updates.
# Use with caution as it provides no end-user control.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "no_end_user_control" {
  display_name                            = "Windows Update Ring - No End User Control"
  description                             = "Automatically install and reboot without end user control"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAndRebootWithoutEndUserControl"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 6: Windows Default

This configuration uses Windows default update behavior, resetting any custom policies.

```terraform
# Scenario 6: Windows Default (Reset)
# This configuration uses the Windows default update behavior, essentially resetting any
# custom update policies to system defaults. Use this to remove custom policies and return
# devices to default Windows Update behavior.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_default" {
  display_name                            = "Windows Update Ring - Windows Default"
  description                             = "Reset to Windows default update behavior"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "windowsDefault"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "disableAllNotifications"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}
```

### Scenario 7: Maximal Assignments

This configuration demonstrates how to assign policies to multiple groups and targets.

```terraform
# Scenario 7: Maximal Assignments
# This configuration demonstrates how to assign a Windows Update Ring to multiple groups
# and built-in targets, including group assignments, all licensed users, all devices,
# and exclusion groups.

# Example groups for assignment (you would use your actual group IDs)
resource "microsoft365_graph_beta_groups_group" "update_ring_group_1" {
  display_name     = "Windows Update Ring - Group 1"
  mail_nickname    = "windows-update-ring-group-1"
  mail_enabled     = false
  security_enabled = true
  description      = "First group for windows update ring assignments"
}

resource "microsoft365_graph_beta_groups_group" "update_ring_group_2" {
  display_name     = "Windows Update Ring - Group 2"
  mail_nickname    = "windows-update-ring-group-2"
  mail_enabled     = false
  security_enabled = true
  description      = "Second group for windows update ring assignments"
}

resource "microsoft365_graph_beta_groups_group" "update_ring_exclusion_group" {
  display_name     = "Windows Update Ring - Exclusion Group"
  mail_nickname    = "windows-update-ring-exclusion-group"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for windows update ring assignments"
}

# Windows Update Ring with comprehensive assignments
resource "microsoft365_graph_beta_device_management_windows_update_ring" "maximal_assignments" {
  display_name                            = "Windows Update Ring - Maximal Assignments"
  description                             = "Demonstrates multiple assignment types"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "notifyDownload"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }

  # Multiple assignment types
  assignments = [
    # Assign to specific group 1
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_group_1.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to specific group 2
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_group_2.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to all licensed users
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to all devices
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Exclude a specific group
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_exclusion_group.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_windows11_upgrade` (Boolean) When TRUE, allows eligible Windows 10 devices to latest Windows 11 release. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.
- `automatic_update_mode` (String) The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime,AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, windowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported.
- `display_name` (String) Admin provided name of the device configuration. Inherited from deviceConfiguration.
- `drivers_excluded` (Boolean) When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.
- `feature_updates_deferral_period_in_days` (Number) Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `feature_updates_rollback_window_in_days` (Number) The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.
- `microsoft_update_service_allowed` (Boolean) When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.
- `quality_updates_deferral_period_in_days` (Number) Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `skip_checks_before_restart` (Boolean) When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.

### Optional

- `active_hours_end` (String) Active Hours End. Part of the Installation Schedule.
- `active_hours_start` (String) Active Hours Start. Part of the Installation Schedule.
- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `business_ready_updates_only` (String) Enable pre-release builds if you want devices to be on a Windows Insider channel.Enabling pre-release builds will cause devices to reboot. Determines which update branch devices will receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, WindowsInsiderBuildSlow, WindowsInsiderBuildRelease.UserDefined equates to 'Not configured' in the gui.all equates to 'Not configured' in the gui.windowsInsiderBuildRelease equates to 'Windows Insider - Release Preview' in the gui.windowsInsiderBuildSlow equates to 'Beta Channel' in the gui.windowsInsiderBuildFast equates to ' Dev Channel' in the gui.
- `deadline_settings` (Attributes) Settings for update installation deadlines and reboot behavior. (see [below for nested schema](#nestedatt--deadline_settings))
- `description` (String) Admin provided description of the Device Configuration. Inherited from deviceConfiguration.
- `feature_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s
- `quality_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `scheduled_install_day` (String) Scheduled Install Day. Possible values are: userDefined, everyday, sunday, monday, tuesday, wednesday, thursday, friday, saturday, noScheduledScan.
- `scheduled_install_time` (String) Scheduled Install Time (in HH:MM:SS format).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `update_notification_level` (String) Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported. Possible values are: notConfigured, defaultNotifications, restartWarningsOnly, disableAllNotifications, unknownFutureValue.
- `update_weeks` (String) Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. Possible values are: userDefined, firstWeek, secondWeek, thirdWeek, fourthWeek, everyWeek, unknownFutureValue.
- `user_pause_access` (String) Specifies whether to enable end user's access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.
- `user_windows_update_scan_access` (String) Specifies whether to disable user's access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. Possible values are: notConfigured, enabled, disabled.

### Read-Only

- `feature_updates_pause_expiry_date_time` (String) The date and time when feature updates pause expires. This value is in ISO 8601 format, in UTC time.
- `feature_updates_pause_start_date` (String) The date when feature updates are paused. This value is in ISO 8601 format, in UTC time.
- `feature_updates_rollback_start_date_time` (String) The date and time when feature updates rollback started. This value is in ISO 8601 format, in UTC time.
- `id` (String) Key of the entity. Inherited from deviceConfiguration.
- `quality_updates_pause_expiry_date_time` (String) The date and time when quality updates pause expires. This value is in ISO 8601 format, in UTC time.
- `quality_updates_pause_start_date` (String) The date when quality updates are paused. This value is in ISO 8601 format, in UTC time.
- `quality_updates_rollback_start_date_time` (String) The date and time when quality updates rollback started. This value is in ISO 8601 format, in UTC time.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--deadline_settings"></a>
### Nested Schema for `deadline_settings`

Required:

- `deadline_for_feature_updates_in_days` (Number) Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_for_quality_updates_in_days` (Number) Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_grace_period_in_days` (Number) Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.
- `postpone_reboot_until_after_deadline` (Boolean) When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.


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