---
page_title: "microsoft365_graph_beta_device_and_app_management_assignment_filter Data Source - microsoft365"
subcategory: "Intune"
description: |-
  
---

# microsoft365_graph_beta_device_and_app_management_assignment_filter (Data Source)

The Microsoft 365 Intune assignment filter data source provides information about a specific Intune assignment filter.

## Example Usage

```terraform
# Basic usage - looking up a single filter by display name
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "by_name" {
  display_name = "Filter | Android Enterprise Device Status Is Rooted"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "windows_vdi" {
  id = "00000000-0000-0000-0000-000000000001"
}

# Example: Create new filter based on existing one (using name lookup)
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "clone_android" {
  display_name = "Clone - Android Rooted Device Filter"
  description  = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.description}"

  # Copy configuration from existing filter
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing all available attributes
output "filter_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.description

    # Filter configuration
    platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
    rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
    assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type

    # Additional metadata
    created_date_time       = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.created_date_time
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.last_modified_date_time
    role_scope_tags         = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.role_scope_tags
  }
}


# Example: Create new filter based on Windows VDI filter (using ID lookup)
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "clone_windows_vdi" {
  display_name = "Clone - Windows VDI Device Filter"
  description  = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"

  # Copy configuration from existing filter
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing Windows VDI filter attributes
output "vdi_filter_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description

    # Filter configuration
    platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
    rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
    assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type

    # Additional metadata
    created_date_time       = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.created_date_time
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.last_modified_date_time
    role_scope_tags         = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags
  }
}


# Use Case 1: Filter Migration - Export multiple filters as JSON for documentation/migration
output "all_filters_export" {
  value = {
    android_filter = {
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.display_name
      config = {
        platform = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
        rule     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
        type     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
      }
    }
    vdi_filter = {
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name
      config = {
        platform = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
        rule     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
        type     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
      }
    }
  }
}

# Use Case 2: Create multiple environment-specific clones with prefix
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "prod_clone" {
  display_name = "PROD - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Production clone of: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"

  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "dev_clone" {
  display_name = "DEV - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Development clone of: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"

  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 3: Create a modified clone with an enhanced rule
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "enhanced_vdi_filter" {
  display_name = "Enhanced - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Enhanced version with additional conditions"

  platform = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  # Original rule with additional conditions
  rule                              = "${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule} and (device.manufacturer -eq \"Microsoft\")"
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 4: Output comparing multiple filters
output "filter_comparison" {
  value = {
    original_vs_enhanced = {
      original_rule = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
      enhanced_rule = "${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule} and (device.manufacturer -eq \"Microsoft\")"
      differences = {
        platform_same = (
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform ==
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
        )
        management_type_same = (
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type ==
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
        )
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name of the assignment filter.
- `id` (String) The unique identifier of the assignment filter.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `assignment_filter_management_type` (String) Indicates filter is applied to either 'devices' or 'apps' management type.
- `created_date_time` (String) The creation time of the assignment filter.
- `description` (String) The description of the assignment filter.
- `last_modified_date_time` (String) Last modified time of the assignment filter.
- `platform` (String) The Intune device management type (platform) for the assignment filter.
- `role_scope_tags` (List of String) Indicates role scope tags assigned for the assignment filter.
- `rule` (String) Rule definition of the assignment filter.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).