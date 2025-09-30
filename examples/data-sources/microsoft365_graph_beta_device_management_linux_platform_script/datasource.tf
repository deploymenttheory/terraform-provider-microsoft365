# Example 1: Get all Linux platform scripts
data "microsoft365_graph_beta_device_management_linux_platform_script" "all_scripts" {
  filter_type = "all"
}

# Example 2: Get a specific Linux platform script by ID
data "microsoft365_graph_beta_device_management_linux_platform_script" "specific_script" {
  filter_type  = "id"
  filter_value = "31fcb6e5-a6a9-4173-8642-5a8572ace9c3" # Replace with actual ID
}

# Example 3: Get Linux platform scripts by display name (partial match)
data "microsoft365_graph_beta_device_management_linux_platform_script" "by_name" {
  filter_type  = "display_name"
  filter_value = "System Config"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_linux_platform_script" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Scripts
output "all_scripts_count" {
  description = "The total number of Linux platform scripts found"
  value       = length(data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items)
}

output "all_scripts_names" {
  description = "List of all Linux platform script names"
  value       = [for script in data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items : script.display_name]
}

output "all_scripts_details" {
  description = "Detailed information for all scripts"
  value = [for script in data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items : {
    id           = script.id
    display_name = script.display_name
    description  = script.description
  }]
}

# Outputs for Specific Script (by ID)
output "specific_script_found" {
  description = "Whether the script with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items) > 0
}

output "specific_script_name" {
  description = "The display name of the script with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items) > 0 ? data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_script_details" {
  description = "Complete details of the script with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Scripts by Name
output "name_filtered_scripts_count" {
  description = "Number of scripts found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items)
}

output "name_filtered_scripts" {
  description = "List of scripts matching the display name filter"
  value = [for script in data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items : {
    id           = script.id
    display_name = script.display_name
    description  = script.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_script" {
  description = "Details of the first script matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first script for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].id
      name        = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_linux_platform_script.specific_script.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_linux_platform_script.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items[0].id
      name        = data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_linux_platform_script.all_scripts.items[0].description
    } : {}
  }
}