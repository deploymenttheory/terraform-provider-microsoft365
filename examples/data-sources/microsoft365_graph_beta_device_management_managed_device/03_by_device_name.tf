# Example: Get managed devices by device name
data "microsoft365_graph_beta_device_management_managed_device" "by_name" {
  device_name = "DESKTOP-WIN-001"
}

# Output: Devices matching the name filter
output "devices_by_name" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_name.items :
    {
      id          = device.id
      device_name = device.device_name
      os          = device.operating_system
      user        = device.user_principal_name
    }
  ]
  description = "List of devices matching the specified device name"
}
