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