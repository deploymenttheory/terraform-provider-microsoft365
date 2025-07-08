---
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages mobile app assignments in Microsoft Intune using the /deviceAppManagement/mobileApps/{mobileAppId}/assignments endpoint. This resource controls how apps are deployed to users and devices, including installation intent (required, available, uninstall), target groups, and platform-specific assignment settings.
---

# microsoft365_graph_beta_device_and_app_management_mobile_app_assignment (Resource)

Manages mobile app assignments in Microsoft Intune using the `/deviceAppManagement/mobileApps/{mobileAppId}/assignments` endpoint. This resource controls how apps are deployed to users and devices, including installation intent (required, available, uninstall), target groups, and platform-specific assignment settings.

## Microsoft Documentation

- [mobileAppAssignment resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta)
- [Create mobileAppAssignment](https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-post-assignments?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementApps.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
########################################################################################
# macOS PKG Assignment Examples
########################################################################################

# Resource for assigning a macos_pkg_app (company_portal) to all licensed users
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_all_users" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allLicensedUsers"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning a macos_pkg_app (company_portal) to all devices
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_all_devices" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allDevices"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning company_portal to a specific group with available install intent
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_group_assignment_available" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "2c39cf3d-78ef-4227-acb1-3a14fc7fbb99"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning company_portal to a specific group with required install intent
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_group_assignment_required" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "7e30b7f0-b2f1-4220-883f-f1d8066eef2d"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

########################################################################################
# Win Get Assignment Examples
########################################################################################

# Resource for assigning a WinGet app (Firefox) to all licensed users
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_all_users" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allLicensedUsers"
    device_and_app_management_assignment_filter_type = "none"
  }

  settings = {
    win_get = {
      notifications = "showAll"
      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2025-06-01T18:00:00Z"
      }
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning a WinGet app (Firefox) to all devices
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_all_devices" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allDevices"
    device_and_app_management_assignment_filter_type = "none"
  }

  settings = {
    win_get = {
      notifications = "showAll"
      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2025-06-01T18:00:00Z"
      }
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with available install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_available" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "2c39cf3d-78ef-4227-acb1-3a14fc7fbb99"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with uninstall install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_uninstall" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "uninstall"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "eadb85bd-6567-4db9-b65c-3f5070d83487"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with required install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_required" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "7e30b7f0-b2f1-4220-883f-f1d8066eef2d"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "exclude"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

########################################################################################
# iOS Store App Assignment Examples
########################################################################################

# Resource for assigning a iOS Store app (Microsoft Edge) to a specific group with required install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "ios_store_app_assignment" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_store_app.example.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
    device_and_app_management_assignment_filter_id   = "471b28c1-8d90-49a2-b639-a47b5f84986d"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    ios_store = {
      is_removable                = true
      prevent_managed_app_backup  = false
      uninstall_on_device_removal = true
      vpn_configuration_id        = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

########################################################################################
# iOS/iPadOS Web Clip Assignment Examples
########################################################################################

# Assignment 1: Available intent to a specific group
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_group_assignment_1" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Assignment 2: Available intent to another group
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_group_assignment_2" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "35d09841-af73-43e6-a59f-024fef1b6b95"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Assignment 3: Available without enrollment intent to a group
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_group_assignment_3" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "availableWithoutEnrollment"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Assignment 4: Available without enrollment with exclusion group
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_group_assignment_4" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "availableWithoutEnrollment"
  source        = "direct"

  target = {
    target_type                                      = "exclusionGroupAssignment"
    group_id                                         = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Assignment 5: Uninstall intent to a group
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_group_assignment_5" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "uninstall"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "e622be02-8c79-48e4-9370-0c78be166eb5"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Assignment 6: Required intent to all licensed users with filter exclusion
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "web_clip_all_users_assignment" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip.company_portal_web_clip.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allLicensedUsers"
    device_and_app_management_assignment_filter_id   = "28b767ca-654c-4605-9371-f1ea044f4207"
    device_and_app_management_assignment_filter_type = "exclude"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `intent` (String) The Intune app install intent defined by the admin. Possible values are:

- **available**: App is available for users to install
- **required**: App is required and will be automatically installed
- **uninstall**: App will be uninstalled
- **availableWithoutEnrollment**: App is available without Intune device enrollment
- `mobile_app_id` (String) The ID of the mobile app to attach the assignment to.
- `source` (String) The resource type which is the source for the assignment. Possible values are: direct, policySets.
- `target` (Attributes) (see [below for nested schema](#nestedatt--target))

### Optional

- `settings` (Attributes) (see [below for nested schema](#nestedatt--settings))
- `source_id` (String) The identifier of the source of the assignment.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of the app assignment associated with the Intune application.

<a id="nestedatt--target"></a>
### Nested Schema for `target`

Required:

- `target_type` (String) The target group type for the application assignment. Possible values are:

- **allDevices**: Target all devices in the tenant
- **allLicensedUsers**: Target all licensed users in the tenant
- **androidFotaDeployment**: Target Android FOTA deployment
- **configurationManagerCollection**: Target System Centre Configuration Manager collection
- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion
- **groupAssignment**: Target a specific Entra ID group

Optional:

- `collection_id` (String) The SCCM group collection ID for the application assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').
- `device_and_app_management_assignment_filter_id` (String) The Id of the scope filter applied to the target assignment.
- `device_and_app_management_assignment_filter_type` (String) The type of scope filter for the target assignment. Defaults to 'none'. Possible values are:

- **include**: Only include devices or users matching the filter
- **exclude**: Exclude devices or users matching the filter
- **none**: No assignment filter applied
- `group_id` (String) The entra ID group ID for the application assignment target. Required when target_type is 'groupAssignment', 'exclusionGroupAssignment', or 'androidFotaDeployment'.


<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `android_managed_store` (Attributes) (see [below for nested schema](#nestedatt--settings--android_managed_store))
- `ios_lob` (Attributes) (see [below for nested schema](#nestedatt--settings--ios_lob))
- `ios_store` (Attributes) (see [below for nested schema](#nestedatt--settings--ios_store))
- `ios_vpp` (Attributes) (see [below for nested schema](#nestedatt--settings--ios_vpp))
- `macos_lob` (Attributes) (see [below for nested schema](#nestedatt--settings--macos_lob))
- `macos_vpp` (Attributes) (see [below for nested schema](#nestedatt--settings--macos_vpp))
- `microsoft_store_for_business` (Attributes) (see [below for nested schema](#nestedatt--settings--microsoft_store_for_business))
- `win32_catalog` (Attributes) (see [below for nested schema](#nestedatt--settings--win32_catalog))
- `win32_lob` (Attributes) (see [below for nested schema](#nestedatt--settings--win32_lob))
- `win_get` (Attributes) (see [below for nested schema](#nestedatt--settings--win_get))
- `windows_app_x` (Attributes) (see [below for nested schema](#nestedatt--settings--windows_app_x))
- `windows_universal_app_x` (Attributes) (see [below for nested schema](#nestedatt--settings--windows_universal_app_x))

<a id="nestedatt--settings--android_managed_store"></a>
### Nested Schema for `settings.android_managed_store`

Optional:

- `android_managed_store_app_track_ids` (List of String) The track IDs to enable for this app assignment.
- `auto_update_mode` (String) The prioritization of automatic updates for this app assignment. Possible values are:

- **default**: Default auto-update mode
- **postponed**: Updates are postponed
- **priority**: Updates are prioritized
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--settings--ios_lob"></a>
### Nested Schema for `settings.ios_lob`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--settings--ios_store"></a>
### Nested Schema for `settings.ios_store`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--settings--ios_vpp"></a>
### Nested Schema for `settings.ios_vpp`

Optional:

- `is_removable` (Boolean) Whether or not the app can be removed by the user. By default, this property is set to FALSE.
- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) Whether or not to uninstall the app when device is removed from Intune. By default, this property is set to FALSE.
- `use_device_licensing` (Boolean) Whether or not to use device licensing. By default, this property is set to FALSE.
- `vpn_configuration_id` (String) The VPN Configuration Id to apply for this app.


<a id="nestedatt--settings--macos_lob"></a>
### Nested Schema for `settings.macos_lob`

Optional:

- `uninstall_on_device_removal` (Boolean) When TRUE, the macOS LOB app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS LOB app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.


<a id="nestedatt--settings--macos_vpp"></a>
### Nested Schema for `settings.macos_vpp`

Optional:

- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, the macOS VPP app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS VPP app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.
- `use_device_licensing` (Boolean) When TRUE indicates that the macOS VPP app should use device-based licensing. When FALSE indicates that the macOS VPP app should use user-based licensing. By default, this property is set to FALSE.


<a id="nestedatt--settings--microsoft_store_for_business"></a>
### Nested Schema for `settings.microsoft_store_for_business`

Optional:

- `use_device_context` (Boolean) When TRUE, indicates that device execution context will be used for the Microsoft Store for Business mobile app. When FALSE, indicates that user context will be used for the Microsoft Store for Business mobile app. By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.


<a id="nestedatt--settings--win32_catalog"></a>
### Nested Schema for `settings.win32_catalog`

Optional:

- `auto_update_settings` (Attributes) The auto-update settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_catalog--auto_update_settings))
- `delivery_optimization_priority` (String) The delivery optimization priority for this app assignment. This setting is not supported in National Cloud environments. Possible values are:

- **notConfigured**: Not configured or background normal delivery optimization priority
- **foreground**: Foreground delivery optimization priority
- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_catalog--install_time_settings))
- `notifications` (String) The notification status for this app assignment. Possible values are:

- **showAll**: Show all notifications
- **showReboot**: Show only reboot notifications
- **hideAll**: Hide all notifications
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_catalog--restart_settings))

<a id="nestedatt--settings--win32_catalog--auto_update_settings"></a>
### Nested Schema for `settings.win32_catalog.auto_update_settings`

Optional:

- `auto_update_superseded_apps_state` (String) The auto-update superseded apps setting for the app assignment. Default value is notConfigured. Possible values are:

- **notConfigured**: Auto-update is not configured
- **enabled**: Auto-update is enabled
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--settings--win32_catalog--install_time_settings"></a>
### Nested Schema for `settings.win32_catalog.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `start_date_time` (String) The time at which the app should be available for installation.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the available and deadline times.


<a id="nestedatt--settings--win32_catalog--restart_settings"></a>
### Nested Schema for `settings.win32_catalog.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--settings--win32_lob"></a>
### Nested Schema for `settings.win32_lob`

Optional:

- `auto_update_settings` (Attributes) The auto-update settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_lob--auto_update_settings))
- `delivery_optimization_priority` (String) The delivery optimization priority for this app assignment. This setting is notsupported in National Cloud environments. Possible values are: notConfigured, foreground.- **notConfigured**: Not configured or background normal delivery optimization priority.
- **foreground**: Foreground delivery optimization priority.
- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_lob--install_time_settings))
- `notifications` (String) The notification status for this app assignment. Possible values are:

- **showAll**: Show all notifications
- **showReboot**: Show only reboot notifications
- **hideAll**: Hide all notifications
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. (see [below for nested schema](#nestedatt--settings--win32_lob--restart_settings))

<a id="nestedatt--settings--win32_lob--auto_update_settings"></a>
### Nested Schema for `settings.win32_lob.auto_update_settings`

Optional:

- `auto_update_superseded_apps_state` (String) The auto-update superseded apps setting for the app assignment. Default value is notConfigured. Possible values are:

- **notConfigured**: Auto-update is not configured
- **enabled**: Auto-update is enabled
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--settings--win32_lob--install_time_settings"></a>
### Nested Schema for `settings.win32_lob.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `start_date_time` (String) The time at which the app should be available for installation.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the available and deadline times.


<a id="nestedatt--settings--win32_lob--restart_settings"></a>
### Nested Schema for `settings.win32_lob.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--settings--win_get"></a>
### Nested Schema for `settings.win_get`

Optional:

- `install_time_settings` (Attributes) (see [below for nested schema](#nestedatt--settings--win_get--install_time_settings))
- `notifications` (String) The notification settings for the assignment. The supported values are 'showAll', 'showReboot', 'hideAll'.
- `restart_settings` (Attributes) (see [below for nested schema](#nestedatt--settings--win_get--restart_settings))

<a id="nestedatt--settings--win_get--install_time_settings"></a>
### Nested Schema for `settings.win_get.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the deadline times.


<a id="nestedatt--settings--win_get--restart_settings"></a>
### Nested Schema for `settings.win_get.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--settings--windows_app_x"></a>
### Nested Schema for `settings.windows_app_x`

Optional:

- `use_device_context` (Boolean) When TRUE, indicates that device execution context will be used for the AppX mobile app. When FALSE, indicates that user context will be used for the AppX mobile app. By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.


<a id="nestedatt--settings--windows_universal_app_x"></a>
### Nested Schema for `settings.windows_universal_app_x`

Optional:

- `use_device_context` (Boolean) If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_win_get_app.example win-get-app-id
```

