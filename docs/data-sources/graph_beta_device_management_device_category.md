---
page_title: "microsoft365_graph_beta_device_management_device_category Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves device categories from Microsoft Intune using the /deviceManagement/deviceCategories endpoint. Device categories help organize devices into logical groups for policy targeting and reporting, enabling users to select categories during enrollment or allowing automatic assignment based on device properties.
---

# microsoft365_graph_beta_device_management_device_category (Data Source)

Retrieves device categories from Microsoft Intune using the `/deviceManagement/deviceCategories` endpoint. Device categories help organize devices into logical groups for policy targeting and reporting, enabling users to select categories during enrollment or allowing automatic assignment based on device properties.

## Microsoft Documentation

- [deviceCategory resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicecategory?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.Read.All`, `DeviceManagementManagedDevices.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Get all device categories
data "microsoft365_graph_beta_device_management_device_category" "all_categories" {
  filter_type = "all"
}

# Example 2: Get a specific device category by ID
data "microsoft365_graph_beta_device_management_device_category" "specific_category" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001" # Replace with actual ID
}

# Example 3: Get device categories by display name (partial match)
data "microsoft365_graph_beta_device_management_device_category" "by_name" {
  filter_type  = "display_name"
  filter_value = "Corporate"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_device_category" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Categories
output "all_categories_count" {
  description = "The total number of device categories found"
  value       = length(data.microsoft365_graph_beta_device_management_device_category.all_categories.items)
}

output "all_categories_names" {
  description = "List of all device category names"
  value       = [for cat in data.microsoft365_graph_beta_device_management_device_category.all_categories.items : cat.display_name]
}

output "all_categories_details" {
  description = "Detailed information for all categories"
  value = [for cat in data.microsoft365_graph_beta_device_management_device_category.all_categories.items : {
    id           = cat.id
    display_name = cat.display_name
    description  = cat.description
  }]
}

# Outputs for Specific Category (by ID)
output "specific_category_found" {
  description = "Whether the category with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_device_category.specific_category.items) > 0
}

output "specific_category_name" {
  description = "The display name of the category with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_device_category.specific_category.items) > 0 ? data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_category_details" {
  description = "Complete details of the category with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_device_category.specific_category.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Categories by Name
output "name_filtered_categories_count" {
  description = "Number of categories found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_device_category.by_name.items)
}

output "name_filtered_categories" {
  description = "List of categories matching the display name filter"
  value = [for cat in data.microsoft365_graph_beta_device_management_device_category.by_name.items : {
    id           = cat.id
    display_name = cat.display_name
    description  = cat.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_category" {
  description = "Details of the first category matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_device_category.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first category for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_device_category.specific_category.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].id
      name        = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_device_category.specific_category.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_device_category.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_device_category.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_device_category.all_categories.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_device_category.all_categories.items[0].id
      name        = data.microsoft365_graph_beta_device_management_device_category.all_categories.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_device_category.all_categories.items[0].description
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

- `items` (Attributes List) The list of Device Categories that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `description` (String) The description of the device category.
- `display_name` (String) The display name of the device category.
- `id` (String) The ID of the device category.