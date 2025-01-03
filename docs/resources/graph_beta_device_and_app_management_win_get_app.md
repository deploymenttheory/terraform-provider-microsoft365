---
page_title: "microsoft365_graph_beta_device_and_app_management_win_get_app Resource - terraform-provider-microsoft365"
subcategory: "Intune"
description: |-
  Manages an Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.
---

# microsoft365_graph_beta_device_and_app_management_win_get_app (Resource)

Manages an Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "whatsapp" {
  package_identifier              = "9NKSQGP7F2NH" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = true

  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "visual_studio_code" {
  package_identifier              = "XP9KHM4BK9FZ7Q" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = true
  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10s"
    update = "10s"
    delete = "10s"
  }

  # App assignments configuration
  assignments = [
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "allDevices"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "allLicensedUsers"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "exclude"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "available"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
        device_and_app_management_assignment_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "uninstall"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "612233b1-55ca-4815-a6b9-5c4aa5a4ac87"
        device_and_app_management_assignment_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        device_and_app_management_assignment_filter_type = "exclude"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `automatically_generate_metadata` (Boolean) When set to `true`, the provider will automatically fetch metadata from the Microsoft Store for Business using the package identifier. This will populate the `display_name`, `description`, `publisher`, and 'icon' fields.
- `install_experience` (Attributes) The install experience settings associated with this application. (see [below for nested schema](#nestedatt--install_experience))
- `package_identifier` (String) The **unique package identifier** for the WinGet/Microsoft Store app from the storefront.

For example, for the app Microsoft Edge Browser URL [https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US](https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US), the package identifier is `xpfftq037jwmhs`.

**Important notes:**
- This identifier is **required** at creation time.
- It **cannot be modified** for existing Terraform-deployed WinGet/Microsoft Store apps.

attempting to modify this value will result in a failed request.

### Optional

- `assignments` (Attributes List) (see [below for nested schema](#nestedatt--assignments))
- `description` (String) A detailed description of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `developer` (String) The developer of the app.
- `display_name` (String) The title of the WinGet app imported from the Microsoft Store for Business.This field value must match the expected title of the app in the Microsoft Store for Business associated with the `package_identifier`.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `information_url` (String) The more information Url.
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin.
- `large_icon` (Attributes) The large icon for the WinGet app, automatically downloaded and set from the Microsoft Store for Business. (see [below for nested schema](#nestedatt--large_icon))
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `publisher` (String) The publisher of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `role_scope_tag_ids` (List of String) List of scope tag ids for this mobile app.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `dependent_app_count` (Number) The total number of dependencies the child app has. This property is read-only.
- `id` (String) The unique graph guid that identifies this resource.Assigned at time of resource creation. This property is read-only.
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group. This property is read-only.
- `last_modified_date_time` (String) The date and time the app was last modified. This property is read-only.
- `manifest_hash` (String) Hash of package metadata properties used to validate that the application matches the metadata in the source repository.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: notPublished, processing, published.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by. This property is read-only.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes. This property is read-only.
- `upload_state` (Number) The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.

<a id="nestedatt--install_experience"></a>
### Nested Schema for `install_experience`

Required:

- `run_as_account` (String) The account type (System or User) that actions should be run as on target devices. Required at creation time.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `intent` (String) The Intune app install intent defined by the admin. Possible values are:

- **available**: App is available for users to install
- **required**: App is required and will be automatically installed
- **uninstall**: App will be uninstalled
- **availableWithoutEnrollment**: App is available without Intune device enrollment
- `source` (String) The resource type which is the source for the assignment. Possible values are: direct, policySets. This property is read-only.
- `target` (Attributes) (see [below for nested schema](#nestedatt--assignments--target))

Optional:

- `settings` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings))

Read-Only:

- `id` (String) The ID of the Intune application associated with this assignment.
- `source_id` (String) The identifier of the source of the assignment. This property is read-only.

<a id="nestedatt--assignments--target"></a>
### Nested Schema for `assignments.target`

Required:

- `target_type` (String) The target group type for the application assignment. Possible values are:

- **allDevices**: Target all devices in the tenant
- **allLicensedUsers**: Target all licensed users in the tenant
- **androidFotaDeployment**: Target Android FOTA deployment
- **configurationManagerCollection**: Target System Centre Configuration Manager collection
- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion
- **groupAssignment**: Target a specific Entra ID group

Optional:

- `collection_id` (String) The SCCM group collection ID for the application assignment target.
- `device_and_app_management_assignment_filter_id` (String) The Id of the scope filter for the target assignment.
- `device_and_app_management_assignment_filter_type` (String) The type of scope filter for the target assignment. Possible values are:

- **include**: Only include devices or users matching the filter
- **exclude**: Exclude devices or users matching the filter
- **none**: No assignment filter applied
- `group_id` (String) The entra ID group ID for the application assignment target.


<a id="nestedatt--assignments--settings"></a>
### Nested Schema for `assignments.settings`

Required:

- `win_get` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--win_get))

Optional:

- `android_managed_store` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--android_managed_store))
- `ios_lob` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--ios_lob))
- `ios_store` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--ios_store))
- `ios_vpp` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--ios_vpp))
- `macos_lob` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--macos_lob))
- `macos_vpp` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--macos_vpp))
- `microsoft_store_for_business` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--microsoft_store_for_business))
- `win32_catalog` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--win32_catalog))
- `win32_lob` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--win32_lob))
- `windows_app_x` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--windows_app_x))
- `windows_universal_app_x` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--windows_universal_app_x))

<a id="nestedatt--assignments--settings--win_get"></a>
### Nested Schema for `assignments.settings.win_get`

Optional:

- `install_time_settings` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--win_get--install_time_settings))
- `notifications` (String) The notification settings for the assignment. Possible values: showAll, showReboot, hideAll
- `restart_settings` (Attributes) (see [below for nested schema](#nestedatt--assignments--settings--win_get--restart_settings))

<a id="nestedatt--assignments--settings--win_get--install_time_settings"></a>
### Nested Schema for `assignments.settings.win_get.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the deadline times.


<a id="nestedatt--assignments--settings--win_get--restart_settings"></a>
### Nested Schema for `assignments.settings.win_get.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--assignments--settings--android_managed_store"></a>
### Nested Schema for `assignments.settings.android_managed_store`

Optional:

- `android_managed_store_app_track_ids` (List of String) The track IDs to enable for this app assignment.
- `auto_update_mode` (String) The prioritization of automatic updates for this app assignment. Possible values are:

- **default**: Default auto-update mode
- **postponed**: Updates are postponed
- **priority**: Updates are prioritized
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--assignments--settings--ios_lob"></a>
### Nested Schema for `assignments.settings.ios_lob`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--assignments--settings--ios_store"></a>
### Nested Schema for `assignments.settings.ios_store`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, this property is set to TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--assignments--settings--ios_vpp"></a>
### Nested Schema for `assignments.settings.ios_vpp`

Optional:

- `is_removable` (Boolean) Whether or not the app can be removed by the user. By default, this property is set to FALSE.
- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to FALSE.
- `uninstall_on_device_removal` (Boolean) Whether or not to uninstall the app when device is removed from Intune. By default, this property is set to FALSE.
- `use_device_licensing` (Boolean) Whether or not to use device licensing. By default, this property is set to FALSE.
- `vpn_configuration_id` (String) The VPN Configuration Id to apply for this app.


<a id="nestedatt--assignments--settings--macos_lob"></a>
### Nested Schema for `assignments.settings.macos_lob`

Optional:

- `uninstall_on_device_removal` (Boolean) When TRUE, the macOS LOB app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS LOB app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.


<a id="nestedatt--assignments--settings--macos_vpp"></a>
### Nested Schema for `assignments.settings.macos_vpp`

Optional:

- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, the macOS VPP app will be uninstalled when the device is removed from Intune management. When FALSE, the macOS VPP app will not be uninstalled when the device is removed from management. By default, this property is set to FALSE.
- `use_device_licensing` (Boolean) When TRUE indicates that the macOS VPP app should use device-based licensing. When FALSE indicates that the macOS VPP app should use user-based licensing. By default, this property is set to FALSE.


<a id="nestedatt--assignments--settings--microsoft_store_for_business"></a>
### Nested Schema for `assignments.settings.microsoft_store_for_business`

Optional:

- `use_device_context` (Boolean) When TRUE, indicates that device execution context will be used for the Microsoft Store for Business mobile app. When FALSE, indicates that user context will be used for the Microsoft Store for Business mobile app. By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.


<a id="nestedatt--assignments--settings--win32_catalog"></a>
### Nested Schema for `assignments.settings.win32_catalog`

Optional:

- `auto_update_settings` (Attributes) The auto-update settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_catalog--auto_update_settings))
- `delivery_optimization_priority` (String) The delivery optimization priority for this app assignment. This setting is not supported in National Cloud environments. Possible values are:

- **notConfigured**: Not configured or background normal delivery optimization priority
- **foreground**: Foreground delivery optimization priority
- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_catalog--install_time_settings))
- `notifications` (String) The notification status for this app assignment. Possible values are:

- **showAll**: Show all notifications
- **showReboot**: Show only reboot notifications
- **hideAll**: Hide all notifications
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_catalog--restart_settings))

<a id="nestedatt--assignments--settings--win32_catalog--auto_update_settings"></a>
### Nested Schema for `assignments.settings.win32_catalog.auto_update_settings`

Optional:

- `auto_update_superseded_apps_state` (String) The auto-update superseded apps setting for the app assignment. Default value is notConfigured. Possible values are:

- **notConfigured**: Auto-update is not configured
- **enabled**: Auto-update is enabled
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--assignments--settings--win32_catalog--install_time_settings"></a>
### Nested Schema for `assignments.settings.win32_catalog.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `start_date_time` (String) The time at which the app should be available for installation.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the available and deadline times.


<a id="nestedatt--assignments--settings--win32_catalog--restart_settings"></a>
### Nested Schema for `assignments.settings.win32_catalog.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--assignments--settings--win32_lob"></a>
### Nested Schema for `assignments.settings.win32_lob`

Optional:

- `auto_update_settings` (Attributes) The auto-update settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--auto_update_settings))
- `delivery_optimization_priority` (String) The delivery optimization priority for this app assignment. This setting is notsupported in National Cloud environments. Possible values are: notConfigured, foreground.- **notConfigured**: Not configured or background normal delivery optimization priority.
- **foreground**: Foreground delivery optimization priority.
- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--install_time_settings))
- `notifications` (String) The notification status for this app assignment. Possible values are:

- **showAll**: Show all notifications
- **showReboot**: Show only reboot notifications
- **hideAll**: Hide all notifications
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--restart_settings))

<a id="nestedatt--assignments--settings--win32_lob--auto_update_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.auto_update_settings`

Optional:

- `auto_update_superseded_apps_state` (String) The auto-update superseded apps setting for the app assignment. Default value is notConfigured. Possible values are:

- **notConfigured**: Auto-update is not configured
- **enabled**: Auto-update is enabled
- **unknownFutureValue**: Reserved for future use


<a id="nestedatt--assignments--settings--win32_lob--install_time_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `start_date_time` (String) The time at which the app should be available for installation.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the available and deadline times.


<a id="nestedatt--assignments--settings--win32_lob--restart_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--assignments--settings--windows_app_x"></a>
### Nested Schema for `assignments.settings.windows_app_x`

Optional:

- `use_device_context` (Boolean) When TRUE, indicates that device execution context will be used for the AppX mobile app. When FALSE, indicates that user context will be used for the AppX mobile app. By default, this property is set to FALSE. Once this property has been set to TRUE it cannot be changed.


<a id="nestedatt--assignments--settings--windows_universal_app_x"></a>
### Nested Schema for `assignments.settings.windows_universal_app_x`

Optional:

- `use_device_context` (Boolean) If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false.




<a id="nestedatt--large_icon"></a>
### Nested Schema for `large_icon`

Optional:

- `type` (String) The MIME type of the app's large icon, automatically populated based on the `package_identifier` when `automatically_generate_metadata` is true. Example: `image/png`
- `value` (String) The icon image to use for the winget app. This field is automatically populated based on the `package_identifier` when `automatically_generate_metadata` is set to true.


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
# Using the provider-default project ID, the import ID is:
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_win_get_app.example win-get-app-id
```

