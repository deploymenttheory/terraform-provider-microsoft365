---
page_title: "microsoft365_graph_beta_device_management_reuseable_policy_setting Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves Reusable Policy Settings from Microsoft Intune with explicit filtering options. Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificate directly to an elevation rule. Like all reusable settings groups for Intune, configurations and changes made to a reusable settings group are automatically passed to the policies that reference the group.
---

# microsoft365_graph_beta_device_management_reuseable_policy_setting (Data Source)

Retrieves Reusable Policy Settings from Microsoft Intune with explicit filtering options. Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificate directly to an elevation rule. Like all reusable settings groups for Intune, configurations and changes made to a reusable settings group are automatically passed to the policies that reference the group.

## Microsoft Documentation

- [deviceManagementReusablePolicySetting resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementreusablepolicysetting?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Get all Reusable Policy Settings
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "all_settings" {
  filter_type = "all"
}

# Example 2: Get a specific Reusable Policy Setting by ID
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "specific_setting" {
  filter_type  = "id"
  filter_value = "64e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Reusable Policy Settings by display name (partial match)
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "by_name" {
  filter_type  = "display_name"
  filter_value = "Certificate"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Settings
output "all_settings_count" {
  description = "The total number of Reusable Policy Settings found"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items)
}

output "all_settings_names" {
  description = "List of all Reusable Policy Setting names"
  value       = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items : setting.display_name]
}

output "all_settings_details" {
  description = "Detailed information for all settings"
  value = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items : {
    id           = setting.id
    display_name = setting.display_name
    description  = setting.description
  }]
}

# Outputs for Specific Setting (by ID)
output "specific_setting_found" {
  description = "Whether the setting with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0
}

output "specific_setting_name" {
  description = "The display name of the setting with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_setting_details" {
  description = "Complete details of the setting with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Settings by Name
output "name_filtered_settings_count" {
  description = "Number of settings found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items)
}

output "name_filtered_settings" {
  description = "List of settings matching the display name filter"
  value = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items : {
    id           = setting.id
    display_name = setting.display_name
    description  = setting.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_setting" {
  description = "Details of the first setting matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first setting for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].description
    } : {}
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of Reusable Policy Settings that match the filter criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `description` (String) The description of the reusable policy setting.
- `display_name` (String) The display name of the reusable policy setting.
- `id` (String) The ID of the reusable policy setting.