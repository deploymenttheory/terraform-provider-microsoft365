---
page_title: "microsoft365_graph_beta_device_and_app_management_application_category Data Source - terraform-provider-microsoft365"
subcategory: "Intune"
description: |-
  Gets information about an Intune mobile app category
---

# microsoft365_graph_beta_device_and_app_management_application_category (Data Source)

The Microsoft 365 Intune application category data source provides information about a specific mobile app category.

## Example Usage

```terraform
# Basic usage - looking up a single application category by display name
data "microsoft365_graph_beta_device_and_app_management_application_category" "by_name" {
  display_name = "Business Apps"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_application_category" "productivity_category" {
  id = "00000000-0000-0000-0000-000000000001"
}

# Example: Create new application category based on existing one (using name lookup)
resource "microsoft365_graph_beta_device_and_app_management_application_category" "clone_business" {
  display_name = "Clone - Business Apps"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing all available attributes
output "category_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.display_name

    # Additional metadata
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.last_modified_date_time
  }
}

# Example: Create new application category based on productivity category (using ID lookup)
resource "microsoft365_graph_beta_device_and_app_management_application_category" "clone_productivity" {
  display_name = "Clone - Productivity Apps"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing productivity category attributes
output "productivity_category_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.display_name

    # Additional metadata
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.last_modified_date_time
  }
}

# Use Case 1: Category Migration - Export multiple categories as JSON for documentation/migration
output "all_categories_export" {
  value = {
    business_category = {
      name = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.display_name
      config = {
        last_modified = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.last_modified_date_time
      }
    }
    productivity_category = {
      name = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.display_name
      config = {
        last_modified = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.last_modified_date_time
      }
    }
  }
}

# Use Case 2: Create multiple environment-specific clones with prefix
resource "microsoft365_graph_beta_device_and_app_management_application_category" "prod_clone" {
  display_name = "PROD - ${data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.display_name}"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_application_category" "dev_clone" {
  display_name = "DEV - ${data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.display_name}"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 3: Create a new standardized category
resource "microsoft365_graph_beta_device_and_app_management_application_category" "standard_apps" {
  display_name = "Standard Applications"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 4: Output comparing multiple categories
output "category_comparison" {
  value = {
    business_vs_productivity = {
      business_name     = data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.display_name
      productivity_name = data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.display_name

      differences = {
        last_modified_same = (
          data.microsoft365_graph_beta_device_and_app_management_application_category.by_name.last_modified_date_time ==
          data.microsoft365_graph_beta_device_and_app_management_application_category.productivity_category.last_modified_date_time
        )
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name of the mobile app category.
- `id` (String) The unique identifier of the mobile app category.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `last_modified_date_time` (String) The date and time when the mobile app category was last modified. This property is read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).