---
page_title: "microsoft365_graph_beta_device_management_windows_update_ring Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves Windows Update Rings from Microsoft Intune with explicit filtering options. Windows Update Rings allow you to define how and when Windows devices receive updates.
---

# microsoft365_graph_beta_device_management_windows_update_ring (Data Source)

Retrieves Windows Update Rings from Microsoft Intune with explicit filtering options. Windows Update Rings allow you to define how and when Windows devices receive updates.

## Microsoft Documentation

- [windowsUpdateForBusinessConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateforbusinessconfiguration?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementConfiguration.ReadWrite.All`

## Example Usage

```terraform
# Example 1: Get all Windows Update Rings
data "microsoft365_graph_beta_device_management_windows_update_ring" "all_rings" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Update Ring by ID
data "microsoft365_graph_beta_device_management_windows_update_ring" "specific_ring" {
  filter_type  = "id"
  filter_value = "a1e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Update Rings by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_update_ring" "by_name" {
  filter_type  = "display_name"
  filter_value = "Pilot"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_update_ring" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Rings
output "all_rings_count" {
  description = "The total number of Windows Update Rings found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items)
}

output "all_rings_names" {
  description = "List of all Windows Update Ring names"
  value       = [for ring in data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items : ring.display_name]
}

output "all_rings_details" {
  description = "Detailed information for all rings"
  value = [for ring in data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items : {
    id           = ring.id
    display_name = ring.display_name
    description  = ring.description
  }]
}

# Outputs for Specific Ring (by ID)
output "specific_ring_found" {
  description = "Whether the ring with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items) > 0
}

output "specific_ring_name" {
  description = "The display name of the ring with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_ring_details" {
  description = "Complete details of the ring with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Rings by Name
output "name_filtered_rings_count" {
  description = "Number of rings found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items)
}

output "name_filtered_rings" {
  description = "List of rings matching the display name filter"
  value = [for ring in data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items : {
    id           = ring.id
    display_name = ring.display_name
    description  = ring.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_ring" {
  description = "Details of the first ring matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first ring for each filtering method
output "update_ring_comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_update_ring.specific_ring.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_update_ring.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items[0].description
    } : {}
  }
}

# Example of using the data in another resource
resource "microsoft365_some_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items) > 0 ? 1 : 0

  name    = "Resource referencing ${data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items[0].display_name}"
  ring_id = data.microsoft365_graph_beta_device_management_windows_update_ring.all_rings.items[0].id

  # Other resource configuration...
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

- `items` (Attributes List) The list of Windows Update Rings that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `description` (String) The description of the Windows Update Ring.
- `display_name` (String) The display name of the Windows Update Ring.
- `id` (String) The ID of the Windows Update Ring.