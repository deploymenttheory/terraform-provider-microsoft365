---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment Resource - terraform-provider-microsoft365"
subcategory: ""
description: |-
  Manages a Mobile App Assignment in Microsoft Intune. Used by different app types to define the assignmentsfor the app. Used by winget_app, windows_web_app, windows_universal_appx, windows_microsoft_edge_app, win32_lob_app, windows_web_app, windows_office_suite_app, managed_ios_store_app, managed_ios_lob_app, managed_android_store_app, managed_ios_lob_app, mac_web_clip, macos_vpp_app, macos_pkg_app, macos_office_suite_app, macos_microsoft_edge_app, macOS_microsoft_defender_app, macOS_lob_app, macOS_dmg_app, ios_vpp_app, ios_store_app, ios_lob_app, ios_ipados_web_clip, android_store_app, android_managed_webstore_app, android_managed_store_app, android_managed_lob_app, android_for_work_app
---

# microsoft365_graph_beta_device_and_app_management_mobile_app_assignment (Resource)

Manages a Mobile App Assignment in Microsoft Intune. Used by different app types to define the assignmentsfor the app. Used by winget_app, windows_web_app, windows_universal_appx, windows_microsoft_edge_app, win32_lob_app, windows_web_app, windows_office_suite_app, managed_ios_store_app, managed_ios_lob_app, managed_android_store_app, managed_ios_lob_app, mac_web_clip, macos_vpp_app, macos_pkg_app, macos_office_suite_app, macos_microsoft_edge_app, macOS_microsoft_defender_app, macOS_lob_app, macOS_dmg_app, ios_vpp_app, ios_store_app, ios_lob_app, ios_ipados_web_clip, android_store_app, android_managed_webstore_app, android_managed_store_app, android_managed_lob_app, android_for_work_app



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `intent` (String) The intent of the assignment. Possible values are: available, required, uninstall, availableWithoutEnrollment.
- `target` (Attributes) The target for this assignment. (see [below for nested schema](#nestedatt--target))

### Optional

- `settings` (Attributes) The settings for this assignment. (see [below for nested schema](#nestedatt--settings))
- `source` (String) The source of the assignment.
- `source_id` (String) The identifier of the source mobile app.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the mobile app assignment.

<a id="nestedatt--target"></a>
### Nested Schema for `target`

Required:

- `type` (String) The type of target. Possible values are: allLicensedUsers, allDevices, group.

Optional:

- `device_and_app_management_assignment_filter_id` (String) The ID of the filter for the target assignment.
- `device_and_app_management_assignment_filter_type` (String) The type of filter for the target assignment. Possible values are: none, include, exclude.
- `group_id` (String) The ID of the group to assign the app to. Required when type is 'group'.


<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `install_time_settings` (Attributes) The install time settings for the assignment. (see [below for nested schema](#nestedatt--settings--install_time_settings))
- `notifications` (String) The notification setting for the assignment. Possible values are: showAll, showReboot, hideAll.
- `restart_settings` (Attributes) The restart settings for the assignment. (see [below for nested schema](#nestedatt--settings--restart_settings))

<a id="nestedatt--settings--install_time_settings"></a>
### Nested Schema for `settings.install_time_settings`

Optional:

- `deadline_date_time` (String) The deadline date and time for the assignment.
- `use_local_time` (Boolean) Indicates whether to use local time for the assignment.


<a id="nestedatt--settings--restart_settings"></a>
### Nested Schema for `settings.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The countdown display before restart in minutes.
- `grace_period_in_minutes` (Number) The grace period before a restart in minutes.
- `restart_notification_snooze_duration_in_minutes` (Number) The snooze duration for the restart notification in minutes.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
