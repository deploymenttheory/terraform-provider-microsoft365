# Example: List all managed devices
data "microsoft365_graph_beta_device_management_managed_device" "all" {
  list_all = true
}

# Output: Total number of managed devices
output "all_managed_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.all.items)
  description = "Total number of managed devices in the tenant"
}

# Output: Device names
output "all_device_names" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all.items :
    device.device_name
  ]
  description = "List of all device names"
}
