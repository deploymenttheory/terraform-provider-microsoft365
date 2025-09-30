# Example 1: Get all Role Scope Tags
data "microsoft365_graph_beta_device_management_role_scope_tag" "all_tags" {
  filter_type = "all"
}

# Example 2: Get a specific Role Scope Tag by ID
data "microsoft365_graph_beta_device_management_role_scope_tag" "specific_tag" {
  filter_type  = "id"
  filter_value = "0da42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Role Scope Tags by display name (partial match)
data "microsoft365_graph_beta_device_management_role_scope_tag" "by_name" {
  filter_type  = "display_name"
  filter_value = "Finance"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_role_scope_tag" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Tags
output "all_tags_count" {
  description = "The total number of Role Scope Tags found"
  value       = length(data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items)
}

output "all_tags_names" {
  description = "List of all Role Scope Tag names"
  value       = [for tag in data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items : tag.display_name]
}

output "all_tags_details" {
  description = "Detailed information for all tags"
  value = [for tag in data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items : {
    id           = tag.id
    display_name = tag.display_name
    description  = tag.description
  }]
}

# Outputs for Specific Tag (by ID)
output "specific_tag_found" {
  description = "Whether the tag with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items) > 0
}

output "specific_tag_name" {
  description = "The display name of the tag with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items) > 0 ? data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_tag_details" {
  description = "Complete details of the tag with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Tags by Name
output "name_filtered_tags_count" {
  description = "Number of tags found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items)
}

output "name_filtered_tags" {
  description = "List of tags matching the display name filter"
  value = [for tag in data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items : {
    id           = tag.id
    display_name = tag.display_name
    description  = tag.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_tag" {
  description = "Details of the first tag matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first tag for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].id
      name        = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_role_scope_tag.specific_tag.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_role_scope_tag.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items[0].id
      name        = data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_role_scope_tag.all_tags.items[0].description
    } : {}
  }
}