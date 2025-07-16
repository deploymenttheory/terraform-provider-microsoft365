# Example: Lookup Windows Driver Update Inventory by ID
data "microsoft365_graph_beta_device_management_windows_driver_update_inventory" "by_id" {
  windows_driver_update_profile_id = "00000000-0000-0000-0000-000000000000" # Replace with your actual profile ID
  id                               = "11111111-1111-1111-1111-111111111111" # Replace with your actual inventory ID

  timeouts {
    read = "5m"
  }
}

# Example: Lookup Windows Driver Update Inventory by Name
data "microsoft365_graph_beta_device_management_windows_driver_update_inventory" "by_name" {
  windows_driver_update_profile_id = "00000000-0000-0000-0000-000000000000" # Replace with your actual profile ID
  name                             = "Intel(R) Wireless Bluetooth(R)"       # Replace with your actual driver name

  timeouts {
    read = "5m"
  }
}

# Output examples
output "driver_inventory_id" {
  value       = data.microsoft365_graph_beta_device_management_windows_driver_update_inventory.by_id.id
  description = "The ID of the driver inventory"
}

output "driver_inventory_name" {
  value       = data.microsoft365_graph_beta_device_management_windows_driver_update_inventory.by_name.name
  description = "The name of the driver inventory"
} 