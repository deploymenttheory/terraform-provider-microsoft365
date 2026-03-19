# Example: Get managed devices by operating system
data "microsoft365_graph_beta_device_management_managed_device" "windows" {
  operating_system = "Windows"
}

# Output: Windows devices
output "windows_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      compliance_state = device.compliance_state
      last_sync        = device.last_sync_date_time
    }
  ]
  description = "List of Windows devices with key information"
}

# Output: Count of Windows devices
output "windows_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.windows.items)
  description = "Total number of Windows devices"
}
