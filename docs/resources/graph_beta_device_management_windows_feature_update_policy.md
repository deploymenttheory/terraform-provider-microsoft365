---
page_title: "microsoft365_graph_beta_device_management_windows_feature_update_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows feature update profiles using the /deviceManagement/windowsFeatureUpdateProfiles endpoint. This resource is used to feature update profiles control major Windows version deployments (like Windows 11 24H2) with rollout scheduling, device eligibility rules, and deployment timing to ensure controlled OS upgrades across managed devices.
---

# microsoft365_graph_beta_device_management_windows_feature_update_policy (Resource)

Manages Windows feature update profiles using the `/deviceManagement/windowsFeatureUpdateProfiles` endpoint. This resource is used to feature update profiles control major Windows version deployments (like Windows 11 24H2) with rollout scheduling, device eligibility rules, and deployment timing to ensure controlled OS upgrades across managed devices.

## Microsoft Documentation

- [windowsFeatureUpdateProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofile?view=graph-rest-beta)
- [Create windowsFeatureUpdateProfile](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-windowsfeatureupdateprofile-create?view=graph-rest-beta)
- [Read windowsFeatureUpdateProfile](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-windowsfeatureupdateprofile-get?view=graph-rest-beta)
- [Update windowsFeatureUpdateProfile](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-windowsfeatureupdateprofile-update?view=graph-rest-beta)
- [Delete windowsFeatureUpdateProfile](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-windowsfeatureupdateprofile-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

### Scenario 1: Make Update Available As Soon As Possible

```terraform
# ==============================================================================
# Make Update Available As Soon As Possible
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update immediately
# without any rollout schedule. The update will be made available to devices
# as soon as possible.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "immediate" {
  display_name                                            = "Windows 11 25H2 - Immediate Deployment"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false
}
```

### Scenario 2: Make Update Available On a Specific Date

```terraform
# ==============================================================================
# Make Update Available On a Specific Date
# ==============================================================================
# This example demonstrates how to schedule a Windows feature update to become
# available on a specific date. The update will not be offered to devices until
# the specified start date.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "scheduled" {
  display_name                                            = "Windows 11 25H2 - Scheduled Deployment"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
  }
}
```

### Scenario 3: Make Update Available Gradually (Phased Rollout)

```terraform
# ==============================================================================
# Make Update Available Gradually (Phased Rollout)
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update gradually
# over time. The update will be offered to devices in stages, starting on the
# offer_start_date and completing by the offer_end_date, with devices being
# offered the update at the specified interval.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "gradual" {
  display_name                                            = "Windows 11 25H2 - Gradual Deployment"
  description                                             = "Phased rollout of Windows 11 25H2 feature update"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
    offer_end_date_time_in_utc   = "2030-01-14T00:00:00Z"
    offer_interval_in_days       = 1
  }

  role_scope_tag_ids = ["0", "1"]
}
```

### Scenario 4: Windows Feature Update Policy with Assignments

```terraform
# ==============================================================================
# Windows Feature Update Policy with Assignments
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update policy with
# group-based assignments. It includes both inclusion and exclusion groups.

# Windows Feature Update Policy with multiple assignments
resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "with_assignments" {
  display_name                                            = "Windows 11 25H2 - With Assignments"
  description                                             = "Feature update deployment with targeted assignments"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the profile.
- `feature_update_version` (String) The feature update version that will be deployed to the devices targeted by this profile. Valid values are: "Windows 11, version 25H2", "Windows 11, version 24H2", "Windows 11, version 23H2". By selecting this Feature update to deploy you are agreeing that when applying this operating system to a device either (1) the applicable Windows license was purchased though volume licensing, or (2) that you are authorized to bind your organization and are accepting on its behalf the relevant Microsoft Software License Terms to be found here https://go.microsoft.com/fwlink/?linkid=2171206.

### Optional

- `assignments` (Attributes Set) Assignments for the Windows Software Update Policies. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `install_feature_updates_optional` (Boolean) If true, the Windows 11 update will become available to users as an optional update. If false, the Windows 11 update will become available to users as a required update
- `install_latest_windows10_on_windows11_ineligible_device` (Boolean) Specifies whether Windows 10 devices that are not eligible for Windows 11 are offered the latest Windows 10 feature updates. Changes to this field require the resource to be replaced.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `rollout_settings` (Attributes) The windows update rollout settings, including offer start date time, offer end date time, and days between each set of offers. (see [below for nested schema](#nestedatt--rollout_settings))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date time that the profile was created.
- `deployable_content_display_name` (String) Friendly display name of the quality update profile deployable content
- `end_of_support_date` (String) The last supported date for a feature update
- `id` (String) The Identifier of the entity.
- `last_modified_date_time` (String) The date time that the profile was last modified.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--rollout_settings"></a>
### Nested Schema for `rollout_settings`

Optional:

- `offer_end_date_time_in_utc` (String) The last group availability of the windows feature update. Must be in RFC3339 format (e.g., '2030-01-13T00:00:00Z').
- `offer_interval_in_days` (Number) The number of days between each set of group offers. The value must be between 1 and 30 and must be equal to or less than the number of days between the 'offer_start_date_time_in_utc' and 'offer_end_date_time_in_utc'.
- `offer_start_date_time_in_utc` (String) The first group availability of the windows feature update. Must be in RFC3339 format (e.g., '2030-01-13T00:00:00Z').


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
terraform import microsoft365_graph_beta_device_and_app_management_windows_feature_update_profile.example windows-feature-update-profile-id
```