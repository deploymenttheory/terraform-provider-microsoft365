# Example: Get managed devices by operating system and version
data "microsoft365_graph_beta_device_management_managed_device" "windows_10" {
  operating_system = "Windows"
  os_version       = "10.0.19045"
}

# Output: Windows 10 devices with specific version
output "windows_10_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows_10.items :
    {
      id          = device.id
      device_name = device.device_name
      os_version  = device.os_version
      compliance  = device.compliance_state
      encrypted   = device.is_encrypted
    }
  ]
  description = "List of Windows 10 devices with the specified version"
}
