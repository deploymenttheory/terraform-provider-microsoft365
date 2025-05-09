# Example 1: Get all Windows Quality Update Expedite Policies
data "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "all_policies" {
  filter_type = "all"
}

# Example 2: Get a specific Windows Quality Update Expedite Policy by ID
data "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "specific_policy" {
  filter_type  = "id"
  filter_value = "c8e42045-31e4-4a8c-95d4-5d7245af782f" # Replace with actual ID
}

# Example 3: Get Windows Quality Update Expedite Policies by display name (partial match)
data "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "by_name" {
  filter_type  = "display_name"
  filter_value = "Critical Updates"
}

# Custom timeout configuration
data "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "with_timeout" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Outputs for All Policies
output "all_expedite_policies_count" {
  description = "The total number of Windows Quality Update Expedite Policies found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items)
}

output "all_expedite_policies_names" {
  description = "List of all Windows Quality Update Expedite Policy names"
  value       = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items : policy.display_name]
}

output "all_expedite_policies_details" {
  description = "Detailed information for all expedite policies"
  value = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items : {
    id           = policy.id
    display_name = policy.display_name
    description  = policy.description
  }]
}

# Outputs for Specific Policy (by ID)
output "specific_expedite_policy_found" {
  description = "Whether the expedite policy with the specified ID was found"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items) > 0
}

output "specific_expedite_policy_name" {
  description = "The display name of the expedite policy with the specified ID"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items) > 0 ? data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].display_name : ""
}

# Using consistent types in conditional
output "specific_expedite_policy_details" {
  description = "Complete details of the expedite policy with the specified ID"
  value = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Outputs for Policies by Name
output "name_filtered_expedite_policies_count" {
  description = "Number of expedite policies found matching the display name filter"
  value       = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items)
}

output "name_filtered_expedite_policies" {
  description = "List of expedite policies matching the display name filter"
  value = [for policy in data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items : {
    id           = policy.id
    display_name = policy.display_name
    description  = policy.description
  }]
}

# Using consistent types in conditional
output "name_filtered_first_expedite_policy" {
  description = "Details of the first expedite policy matching the display name filter (if any)"
  value = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items) > 0 ? {
    id           = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].id
    display_name = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].display_name
    description  = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].description
    found        = true
    } : {
    id           = ""
    display_name = ""
    description  = ""
    found        = false
  }
}

# Simple output showing the first policy for each filtering method
output "expedite_policy_comparison_summary" {
  description = "Summary comparison of results from each filtering method"
  value = {
    by_id = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.specific_policy.items[0].description
    } : {}

    by_name = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.by_name.items[0].description
    } : {}

    all_first = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items) > 0 ? {
      id          = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items[0].id
      name        = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items[0].display_name
      description = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items[0].description
    } : {}
  }
}

# Example of using the data in another resource
resource "microsoft365_some_other_resource" "example" {
  count = length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items) > 0 ? 1 : 0

  name               = "Resource referencing ${data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items[0].display_name}"
  expedite_policy_id = data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items[0].id

  # Other resource configuration...
}

# Example of using multiple policies in a loop
resource "microsoft365_some_policy_association" "associations" {
  for_each = {
    for idx, policy in data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items : policy.id => policy
    if length(data.microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.all_policies.items) > 0
  }

  policy_id   = each.key
  policy_name = each.value.display_name
  enabled     = true

  # Other association configuration...
}