# Example 1: Get all Windows Feature Update Profiles
data "microsoft365_graph_beta_device_management_windows_feature_update_policy" "all_profiles" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Feature Update Profile by ID
data "microsoft365_graph_beta_device_management_windows_feature_update_policy" "specific_profile" {
  filter_type  = "id"
  filter_value = "b9e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Feature Update Profiles by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_feature_update_policy" "by_name" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_feature_update_policy" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Profiles
output "all_profiles_count" {
  description = "The total number of Windows Feature Update Profiles found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items)
}

output "all_profiles_names" {
  description = "List of all Windows Feature Update Profile names"
  value       = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items : profile.display_name]
}

output "all_profiles_details" {
  description = "Detailed information for all profiles"
  value = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items : {
    id           = profile.id
    display_name = profile.display_name
    description  = profile.description
  }]
}

# Outputs for Specific Profile (by ID)
output "specific_profile_found" {
  description = "Whether the profile with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items) > 0
}

output "specific_profile_name" {
  description = "The display name of the profile with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_profile_details" {
  description = "Complete details of the profile with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Profiles by Name
output "name_filtered_profiles_count" {
  description = "Number of profiles found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items)
}

output "name_filtered_profiles" {
  description = "List of profiles matching the display name filter"
  value = [for profile in data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items : {
    id           = profile.id
    display_name = profile.display_name
    description  = profile.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_profile" {
  description = "Details of the first profile matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first profile for each filtering method
output "comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.specific_profile.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.all_profiles.items[0].description
    } : {}
  }
}

# Using the data in a configuration
resource "microsoft365_some_other_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items) > 0 ? 1 : 0

  name       = "Resource using ${data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].display_name}"
  profile_id = data.microsoft365_graph_beta_device_management_windows_feature_update_policy.by_name.items[0].id

  # Other resource configuration...
}