# Example: Get managed devices using OData filter with AND logic
data "microsoft365_graph_beta_device_management_managed_device" "compliant_windows" {
  odata_query = "operatingSystem eq 'Windows' and complianceState eq 'compliant'"
}

# Output: Compliant Windows devices
output "compliant_windows_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.compliant_windows.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      compliance_state = device.compliance_state
      is_encrypted     = device.is_encrypted
    }
  ]
  description = "List of compliant Windows devices"
}
