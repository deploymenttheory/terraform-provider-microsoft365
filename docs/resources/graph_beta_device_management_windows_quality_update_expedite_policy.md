---
page_title: "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows quality update expedite policies in Microsoft Intune using the /deviceManagement/windowsQualityUpdateProfiles endpoint. This resource is used to enable accelerated deployment of critical Windows quality updates with forced reboot enforcement for urgent security patches.
---

# microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy (Resource)

Manages Windows quality update expedite policies in Microsoft Intune using the `/deviceManagement/windowsQualityUpdateProfiles` endpoint. This resource is used to enable accelerated deployment of critical Windows quality updates with forced reboot enforcement for urgent security patches.

## Microsoft Documentation

- [expediteWindowsQualityUpdateSettings resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-expeditewindowsqualityupdatesettings?view=graph-rest-beta)
- [Create expediteWindowsQualityUpdateSettings](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-expeditewindowsqualityupdatesettings-create?view=graph-rest-beta)
- [Read expediteWindowsQualityUpdateSettings](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-expeditewindowsqualityupdatesettings-get?view=graph-rest-beta)
- [Update expediteWindowsQualityUpdateSettings](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-expeditewindowsqualityupdatesettings-update?view=graph-rest-beta)
- [Delete expediteWindowsQualityUpdateSettings](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-expeditewindowsqualityupdatesettings-delete?view=graph-rest-beta)

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

### Maximal Configuration (No Assignments)

```terraform
# Example: Windows Quality Update Expedite Policy - Maximal Configuration (No Assignments)
resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "maximal_example" {
  display_name       = "Critical Security Update Expedite Policy"
  description        = "Expedited deployment for critical security updates - January 2025"
  role_scope_tag_ids = ["0", "1"]

  # Required: Expedited update settings
  # Defines which quality update to expedite and reboot behavior
  expedited_update_settings = {
    # Quality update release to expedite
    # Valid values: "2025-12-09T00:00:00Z", "2025-11-20T00:00:00Z"
    quality_update_release = "2025-12-09T00:00:00Z"

    # If a reboot is required, select the number of days before it's enforced
    # Valid values: 0, 1, or 2
    # 0 = Immediate reboot after installation
    # 1 = Allow 1 day for user-initiated reboot
    # 2 = Allow 2 days for user-initiated reboot
    days_until_forced_reboot = 1
  }

  # Optional: Custom timeouts for resource operations
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}
```

### Maximal Configuration with Assignments

```terraform
# Example: Windows Quality Update Expedite Policy - Maximal Configuration with Assignments
resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "maximal_with_assignments" {
  display_name       = "Production Critical Update Expedite Policy"
  description        = "Expedited deployment for critical security updates targeting production devices with exclusions for test environments"
  role_scope_tag_ids = ["0", "1"]

  # Required: Expedited update settings
  expedited_update_settings = {
    # Quality update release to expedite
    # Valid values: "2025-12-09T00:00:00Z", "2025-11-20T00:00:00Z"
    quality_update_release = "2025-11-20T00:00:00Z"

    # Force reboot after 2 days to ensure update completion
    days_until_forced_reboot = 2
  }

  # Assignments: Target specific groups and exclude others
  assignments = [
    # Primary production group assignment
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111" # Production Devices Group
    },
    # Secondary production group assignment
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222" # Executive Devices Group
    },
    # Exclude test environment group
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333" # Test Devices Group
    },
    # Exclude development environment group
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444" # Development Devices Group
    }
  ]

  # Optional: Custom timeouts for resource operations
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the Windows Quality Update Profile (expedite policy).
- `expedited_update_settings` (Attributes) Expedited Quality update settings. (see [below for nested schema](#nestedatt--expedited_update_settings))

### Optional

- `assignments` (Attributes Set) Assignments for the Windows Software Update Policies. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date time that the profile was created.
- `deployable_content_display_name` (String) Friendly display name of the quality update profile deployable content
- `id` (String) The Intune Windows Quality Update Profile (expedite policy) profile id.
- `last_modified_date_time` (String) The date time that the profile was last modified.
- `release_date_display_name` (String) Friendly release date to display for a Quality Update release

<a id="nestedatt--expedited_update_settings"></a>
### Nested Schema for `expedited_update_settings`

Required:

- `days_until_forced_reboot` (Number) If a reboot is required, select the number of days before it's enforced. Valid values are: 0, 1, and 2.
- `quality_update_release` (String) Expedite installation of quality updates if device OS version less than the quality update release identifier. Value must be a valid ISO 8601 datetime format (e.g., 2025-12-09T00:00:00Z). Valid values as of December 2025: 2025-12-09T00:00:00Z, 2025-11-20T00:00:00Z


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


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
terraform import microsoft365_graph_beta_device_and_app_management_windows_quality_update_expedite_policy.example 00000000-0000-0000-0000-000000000001
```