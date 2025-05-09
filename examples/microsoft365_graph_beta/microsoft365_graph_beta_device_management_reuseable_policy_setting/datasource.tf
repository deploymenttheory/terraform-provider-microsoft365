# Example 1: Get all Reusable Policy Settings
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "all_settings" {
  filter_type = "all"
}

# Example 2: Get a specific Reusable Policy Setting by ID
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "specific_setting" {
  filter_type  = "id"
  filter_value = "64e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Reusable Policy Settings by display name (partial match)
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "by_name" {
  filter_type  = "display_name"
  filter_value = "Certificate"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_reuseable_policy_setting" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Settings
output "all_settings_count" {
  description = "The total number of Reusable Policy Settings found"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items)
}

output "all_settings_names" {
  description = "List of all Reusable Policy Setting names"
  value       = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items : setting.display_name]
}

output "all_settings_details" {
  description = "Detailed information for all settings"
  value = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items : {
    id           = setting.id
    display_name = setting.display_name
    description  = setting.description
  }]
}

# Outputs for Specific Setting (by ID)
output "specific_setting_found" {
  description = "Whether the setting with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0
}

output "specific_setting_name" {
  description = "The display name of the setting with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_setting_details" {
  description = "Complete details of the setting with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Settings by Name
output "name_filtered_settings_count" {
  description = "Number of settings found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items)
}

output "name_filtered_settings" {
  description = "List of settings matching the display name filter"
  value = [for setting in data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items : {
    id           = setting.id
    display_name = setting.display_name
    description  = setting.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_setting" {
  description = "Details of the first setting matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first setting for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.specific_setting.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].id
      name        = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_reuseable_policy_setting.all_settings.items[0].description
    } : {}
  }
}