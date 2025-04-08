# Basic usage - looking up a single device category by display name
data "microsoft365_graph_beta_device_and_app_management_device_category" "by_name" {
  display_name = "Corporate Laptops"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_device_category" "byod_category" {
  id = "00000000-0000-0000-0000-000000000001"
}

# Example: Create new device category based on existing one (using name lookup)
resource "microsoft365_graph_beta_device_and_app_management_device_category" "clone_corporate" {
  display_name       = "Clone - Corporate Laptops"
  description        = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.description}"
  role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.role_scope_tag_ids

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
    id           = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.description

    # Additional metadata
    role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.role_scope_tag_ids
  }
}

# Example: Create new device category based on BYOD category (using ID lookup)
resource "microsoft365_graph_beta_device_and_app_management_device_category" "clone_byod" {
  display_name       = "Clone - BYOD Devices"
  description        = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description}"
  role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing BYOD category attributes
output "byod_category_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description

    # Additional metadata
    role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids
  }
}

# Use Case 1: Category Migration - Export multiple categories as JSON for documentation/migration
output "all_categories_export" {
  value = {
    corporate_category = {
      name = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.display_name
      config = {
        description = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.description
        tags        = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.role_scope_tag_ids
      }
    }
    byod_category = {
      name = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.display_name
      config = {
        description = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description
        tags        = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids
      }
    }
  }
}

# Use Case 2: Create multiple environment-specific clones with prefix
resource "microsoft365_graph_beta_device_and_app_management_device_category" "prod_clone" {
  display_name       = "PROD - ${data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.display_name}"
  description        = "Production clone of: ${data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description}"
  role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_device_category" "dev_clone" {
  display_name       = "DEV - ${data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.display_name}"
  description        = "Development clone of: ${data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description}"
  role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 3: Create a new standardized category
resource "microsoft365_graph_beta_device_and_app_management_device_category" "standard_workstation" {
  display_name       = "Standard Workstation"
  description        = "Standard corporate workstation for office employees"
  role_scope_tag_ids = ["0", "9"] # Example role scope tags

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
    corporate_vs_byod = {
      corporate_name = data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.display_name
      byod_name      = data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.display_name

      differences = {
        description_same = (
          data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.description ==
          data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.description
        )
        tags_count_corporate = length(data.microsoft365_graph_beta_device_and_app_management_device_category.by_name.role_scope_tag_ids)
        tags_count_byod      = length(data.microsoft365_graph_beta_device_and_app_management_device_category.byod_category.role_scope_tag_ids)
      }
    }
  }
}