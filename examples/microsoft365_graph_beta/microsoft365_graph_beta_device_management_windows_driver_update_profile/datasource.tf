# Example: Lookup Windows Driver Update Profile by ID
data "microsoft365_graph_beta_device_management_windows_driver_update_profile" "by_id" {
  id = "00000000-0000-0000-0000-000000000000" # Replace with your actual profile ID

  timeouts {
    read = "5m"
  }
}

# Example: Lookup Windows Driver Update Profile by Display Name
data "microsoft365_graph_beta_device_management_windows_driver_update_profile" "by_name" {
  display_name = "Windows 11 Driver Updates" # Replace with your actual profile name

  timeouts {
    read = "5m"
  }
}

# Output examples
output "profile_id" {
  value       = data.microsoft365_graph_beta_device_management_windows_driver_update_profile.by_id.id
  description = "The ID of the Windows driver update profile"
}

output "profile_description" {
  value       = data.microsoft365_graph_beta_device_management_windows_driver_update_profile.by_name.description
  description = "The description of the Windows driver update profile"
}

output "profile_role_scope_tag_ids" {
  value       = data.microsoft365_graph_beta_device_management_windows_driver_update_profile.by_name.role_scope_tag_ids
  description = "The role scope tag IDs associated with the profile"
}

# Example: Use the profile data source to look up driver inventories
data "microsoft365_graph_beta_device_management_windows_driver_update_inventory" "example" {
  windows_driver_update_profile_id = data.microsoft365_graph_beta_device_management_windows_driver_update_profile.by_name.id
  name                             = "Intel(R) Wireless Bluetooth(R)" # Replace with your actual driver name
}
