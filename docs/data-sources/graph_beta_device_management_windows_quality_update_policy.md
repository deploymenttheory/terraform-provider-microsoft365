---
page_title: "microsoft365_graph_beta_device_management_windows_quality_update_policy Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves information about a Windows Driver Update Profile in Microsoft Intune.
---

# microsoft365_graph_beta_device_management_windows_quality_update_policy (Data Source)

Retrieves information about a Windows Driver Update Profile in Microsoft Intune.

## Microsoft Documentation

- [windowsQualityUpdatePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsqualityupdatepolicy?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementConfiguration.ReadWrite.All`

## Example Usage

```terraform
# Example 1: Get all Windows Quality Update Policies
data "microsoft365_graph_beta_device_management_windows_quality_update_policy" "all_policies" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Quality Update Policy by ID
data "microsoft365_graph_beta_device_management_windows_quality_update_policy" "specific_policy" {
  filter_type  = "id"
  filter_value = "a7e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Quality Update Policies by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_quality_update_policy" "by_name" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_quality_update_policy" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Policies
output "all_policies_count" {
  description = "The total number of Windows Quality Update Policies found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items)
}

output "all_policies_names" {
  description = "List of all Windows Quality Update Policy names"
  value       = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items : policy.display_name]
}

output "all_policies_details" {
  description = "Detailed information for all policies"
  value = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items : {
    id           = policy.id
    display_name = policy.display_name
    description  = policy.description
  }]
}

# Outputs for Specific Policy (by ID)
output "specific_policy_found" {
  description = "Whether the policy with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items) > 0
}

output "specific_policy_name" {
  description = "The display name of the policy with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_policy_details" {
  description = "Complete details of the policy with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Policies by Name
output "name_filtered_policies_count" {
  description = "Number of policies found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items)
}

output "name_filtered_policies" {
  description = "List of policies matching the display name filter"
  value = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items : {
    id           = policy.id
    display_name = policy.display_name
    description  = policy.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_policy" {
  description = "Details of the first policy matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first policy for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.specific_policy.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_policy.all_policies.items[0].description
    } : {}
  }
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
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Driver Update entity.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).