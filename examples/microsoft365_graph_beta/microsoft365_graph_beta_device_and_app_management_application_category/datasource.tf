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