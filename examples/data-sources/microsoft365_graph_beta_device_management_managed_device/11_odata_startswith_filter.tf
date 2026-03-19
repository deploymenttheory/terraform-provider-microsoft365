# Example: Get managed devices using OData startswith filter
data "microsoft365_graph_beta_device_management_managed_device" "desktop_devices" {
  odata_query = "startswith(deviceName,'DESKTOP')"
}

# Output: Desktop devices
output "desktop_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.desktop_devices.items :
    {
      id          = device.id
      device_name = device.device_name
      os          = device.operating_system
      user        = device.user_principal_name
    }
  ]
  description = "List of devices with names starting with 'DESKTOP'"
}
