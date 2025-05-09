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