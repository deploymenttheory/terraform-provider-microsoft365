# Example 1: Get all assignment filters
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "all_filters" {
  filter_type = "all"
}

# Example 2: Get a specific assignment filter by ID
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "specific_filter" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001" # Replace with actual ID
}

# Example 3: Get assignment filters by display name (partial match)
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "by_name" {
  filter_type  = "display_name"
  filter_value = "Purpose-built Specialty Devices On Android Device Administrator"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "with_timeout" {
  filter_type = "all"
  
  timeouts = {
    read = "1m" 
  }
}

# Outputs for All Filters
output "all_filters_count" {
  description = "The total number of assignment filters found"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items)
}

output "all_filters_names" {
  description = "List of all assignment filter names"
  value       = [for filter in data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items : filter.display_name]
}

output "all_filters_details" {
  description = "Detailed information for all filters"
  value       = [for filter in data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items : {
    id          = filter.id
    display_name = filter.display_name
    description = filter.description
  }]
}

# Outputs for Specific Filter (by ID)
output "specific_filter_found" {
  description = "Whether the filter with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items) > 0
}

output "specific_filter_name" {
  description = "The display name of the filter with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_filter_details" {
  description = "Complete details of the filter with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].description
    found        = true
  } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Filters by Name
output "name_filtered_filters_count" {
  description = "Number of filters found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items)
}

output "name_filtered_filters" {
  description = "List of filters matching the display name filter"
  value       = [for filter in data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items : {
    id           = filter.id
    display_name = filter.display_name
    description  = filter.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_filter" {
  description = "Details of the first filter matching the display name filter (if any)"
  value       = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].description
    found        = true
  } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first filter for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items) > 0 ? {
      id = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].display_name
      description = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.specific_filter.items[0].description
    } : {}
    
    by_name = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items) > 0 ? {
      id = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.items[0].description
    } : {}
    
    all_first = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items) > 0 ? {
      id = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items[0].id
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items[0].display_name
      description = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items[0].description
    } : {}
  }
}

# Using assignment filter data in another resource
resource "microsoft365_some_other_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items) > 0 ? 1 : 0
  
  name = "Resource using ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items[0].display_name}"
  filter_id = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.all_filters.items[0].id
  
  # Other resource configuration...
}