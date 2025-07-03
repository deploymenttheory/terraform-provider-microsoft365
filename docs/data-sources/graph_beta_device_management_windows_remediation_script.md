---
page_title: "microsoft365_graph_beta_device_management_windows_remediation_script Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves Windows Remediation Scripts from Microsoft Intune with explicit filtering options. Windows Remediation Scripts are PowerShell scripts that can be deployed to devices to help remediate issues.
---

# microsoft365_graph_beta_device_management_windows_remediation_script (Data Source)

The Microsoft 365 Intune windows update catalog item data source provides information about a windows updates. Can be filtered by
id, display_name , end_of_support_date and release_date_time.

## Example Usage

```terraform
# Example 1: Get all Windows Remediation Scripts
data "microsoft365_graph_beta_device_management_windows_remediation_script" "all_scripts" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Remediation Script by ID
data "microsoft365_graph_beta_device_management_windows_remediation_script" "specific_script" {
  filter_type  = "id"
  filter_value = "a5e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Remediation Scripts by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_remediation_script" "by_name" {
  filter_type  = "display_name"
  filter_value = "Fix Windows"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_remediation_script" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Scripts
output "all_remediation_scripts_count" {
  description = "The total number of Windows Remediation Scripts found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items)
}

output "all_remediation_scripts_names" {
  description = "List of all Windows Remediation Script names"
  value       = [for script in data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items : script.display_name]
}

output "all_remediation_scripts_details" {
  description = "Detailed information for all remediation scripts"
  value = [for script in data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items : {
    id           = script.id
    display_name = script.display_name
    description  = script.description
  }]
}

# Outputs for Specific Script (by ID)
output "specific_remediation_script_found" {
  description = "Whether the script with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items) > 0
}

output "specific_remediation_script_name" {
  description = "The display name of the script with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_remediation_script_details" {
  description = "Complete details of the script with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Scripts by Name
output "name_filtered_remediation_scripts_count" {
  description = "Number of scripts found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items)
}

output "name_filtered_remediation_scripts" {
  description = "List of scripts matching the display name filter"
  value = [for script in data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items : {
    id           = script.id
    display_name = script.display_name
    description  = script.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_remediation_script" {
  description = "Details of the first script matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first script for each filtering method
output "remediation_script_comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_remediation_script.specific_script.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_remediation_script.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items[0].description
    } : {}
  }
}

# Example of using the data in another resource
resource "microsoft365_some_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items) > 0 ? 1 : 0

  name      = "Resource referencing ${data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items[0].display_name}"
  script_id = data.microsoft365_graph_beta_device_management_windows_remediation_script.all_scripts.items[0].id

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

- `items` (Attributes List) The list of Windows Remediation Scripts that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `description` (String) The description of the Windows Remediation Script.
- `display_name` (String) The display name of the Windows Remediation Script.
- `id` (String) The ID of the Windows Remediation Script.