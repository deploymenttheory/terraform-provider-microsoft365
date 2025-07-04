---
page_title: "microsoft365_graph_beta_device_and_app_management_application_category Data Source - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Retrieves Application Categories from Microsoft Intune with explicit filtering options.
---

# microsoft365_graph_beta_device_and_app_management_application_category (Data Source)

Retrieves Application Categories from Microsoft Intune with explicit filtering options.

## Microsoft Documentation

- [mobileAppCategory resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementApps.Read.All`, `DeviceManagementApps.ReadWrite.All`

## Example Usage

```terraform
# Example 1: Get all application categories
data "microsoft365_graph_beta_device_and_app_management_application_category" "all_categories" {
  filter_type = "all"
}

# Example 2: Get a specific application category by ID
data "microsoft365_graph_beta_device_and_app_management_application_category" "specific_category" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001" # Replace with actual ID
}

# Example 3: Get application categories by display name (partial match)
data "microsoft365_graph_beta_device_and_app_management_application_category" "by_name" {
  filter_type  = "display_name"
  filter_value = "Computer Management"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_and_app_management_application_category" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Categories
output "all_categories_count" {
  description = "The total number of application categories found"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items)
}

output "all_categories_names" {
  description = "List of all application category names"
  value       = [for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items : cat.display_name]
}

output "all_categories_details" {
  description = "Detailed information for all categories"
  value = [for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items : {
    id            = cat.id
    display_name  = cat.display_name
    last_modified = cat.last_modified_date_time
  }]
}

# Outputs for Specific Category (by ID)
output "specific_category_found" {
  description = "Whether the category with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items) > 0
}

output "specific_category_name" {
  description = "The display name of the category with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].display_name : ""
}

# Use consistent types in conditional
output "specific_category_details" {
  description = "Complete details of the category with the specified ID"
  value = length(data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items) > 0 ? {
    id            = data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].id
    display_name  = data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].display_name
    last_modified = data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].last_modified_date_time
    found         = true
    } : {
    id            = ""
    display_name  = ""
    last_modified = ""
    found         = false
  }
}

# Outputs for Categories by Name
output "name_filtered_categories_count" {
  description = "Number of categories found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items)
}

output "name_filtered_categories" {
  description = "List of categories matching the display name filter"
  value = [for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items : {
    id            = cat.id
    display_name  = cat.display_name
    last_modified = cat.last_modified_date_time
  }]
}

# Use consistent types in conditional
output "name_filtered_first_category" {
  description = "Details of the first category matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items) > 0 ? {
    id            = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items[0].id
    display_name  = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items[0].display_name
    last_modified = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items[0].last_modified_date_time
    found         = true
    } : {
    id            = ""
    display_name  = ""
    last_modified = ""
    found         = false
  }
}

# Example of using the data in conditional outputs
output "computer_management_category_id" {
  description = "ID of the Computer Management category, if found"
  value       = length([for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items : cat if cat.display_name == "Computer Management"]) > 0 ? [for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items : cat.id if cat.display_name == "Computer Management"][0] : ""
}

# Demonstration of comparing results between different filtering methods
output "filtering_comparison" {
  description = "Verifies that filtering by ID and name return expected results"
  value = {
    id_filter_matched_name  = length(data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items) > 0 && length(data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items) > 0 ? contains([for cat in data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items : cat.id], data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].id) : false
    name_filter_works       = length(data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items) > 0
    categories_with_timeout = length(data.microsoft365_graph_beta_device_and_app_management_application_category.with_timeout.items)
  }
}

# Simple output showing the first category for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items) > 0 ? {
      id   = data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_application_category.specific_category.items[0].display_name
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items) > 0 ? {
      id   = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.items[0].display_name
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items) > 0 ? {
      id   = data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_application_category.all_categories.items[0].display_name
    } : {}
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `last_modified_date_time`.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'. For date filters, use RFC3339 format (e.g., '2023-01-01T00:00:00Z').
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of Application Categories that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `display_name` (String) The display name of the application category.
- `id` (String) The ID of the application category.
- `last_modified_date_time` (String) The date and time when the application category was last modified.