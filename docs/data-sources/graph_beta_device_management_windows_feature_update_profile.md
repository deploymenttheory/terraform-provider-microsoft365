---
page_title: "microsoft365_graph_beta_device_management_windows_feature_update_profile Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves information about a Windows Feature Update Profile in Microsoft Intune.
---

# microsoft365_graph_beta_device_management_windows_feature_update_profile (Data Source)

Retrieves information about a Windows Feature Update Profile in Microsoft Intune.

## Microsoft Documentation

- [windowsFeatureUpdateProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofile?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementConfiguration.ReadWrite.All`

## Example Usage

```terraform
# Example 1: Get all Windows Feature Update Profiles
data "microsoft365_graph_beta_device_management_windows_feature_update_profile" "all_profiles" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Feature Update Profile by ID
data "microsoft365_graph_beta_device_management_windows_feature_update_profile" "specific_profile" {
  filter_type  = "id"
  filter_value = "b9e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Feature Update Profiles by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_feature_update_profile" "by_name" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_feature_update_profile" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Profiles
output "all_profiles_count" {
  description = "The total number of Windows Feature Update Profiles found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items)
}

output "all_profiles_names" {
  description = "List of all Windows Feature Update Profile names"
  value       = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items : profile.display_name]
}

output "all_profiles_details" {
  description = "Detailed information for all profiles"
  value = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items : {
    id           = profile.id
    display_name = profile.display_name
    description  = profile.description
  }]
}

# Outputs for Specific Profile (by ID)
output "specific_profile_found" {
  description = "Whether the profile with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items) > 0
}

output "specific_profile_name" {
  description = "The display name of the profile with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_profile_details" {
  description = "Complete details of the profile with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Profiles by Name
output "name_filtered_profiles_count" {
  description = "Number of profiles found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items)
}

output "name_filtered_profiles" {
  description = "List of profiles matching the display name filter"
  value = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items : {
    id           = profile.id
    display_name = profile.display_name
    description  = profile.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_profile" {
  description = "Details of the first profile matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first profile for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_profile.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.all_profiles.items[0].description
    } : {}
  }
}

# Using the data in a configuration
resource "microsoft365_some_other_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items) > 0 ? 1 : 0

  name       = "Resource using ${data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].display_name}"
  profile_id = data.microsoft365_graph_beta_device_management_windows_feature_update_profile.by_name.items[0].id

  # Other resource configuration...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name for the profile.
- `id` (String) The ID of the profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The description of the profile which is specified by the user.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Feature Update entity.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).