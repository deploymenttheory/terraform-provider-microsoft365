# Example: Get managed devices by serial number
data "microsoft365_graph_beta_device_management_managed_device" "by_serial" {
  serial_number = "SN-WIN-001"
}

# Output: Devices with matching serial number
output "devices_by_serial" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_serial.items :
    {
      id            = device.id
      device_name   = device.device_name
      serial_number = device.serial_number
      manufacturer  = device.manufacturer
      model         = device.model
      user          = device.user_display_name
    }
  ]
  description = "List of devices matching the specified serial number"
}

# Note: Multiple devices may share the same serial number in some scenarios
output "serial_number_device_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.by_serial.items)
  description = "Number of devices with this serial number"
}
