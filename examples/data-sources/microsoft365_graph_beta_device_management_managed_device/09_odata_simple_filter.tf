# Example: Get managed devices using a simple OData filter
data "microsoft365_graph_beta_device_management_managed_device" "compliant_devices" {
  odata_query = "complianceState eq 'compliant'"
}

# Output: Compliant devices
output "compliant_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.compliant_devices.items :
    {
      id               = device.id
      device_name      = device.device_name
      compliance_state = device.compliance_state
      operating_system = device.operating_system
    }
  ]
  description = "List of compliant devices"
}

# Output: Compliant device count
output "compliant_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.compliant_devices.items)
  description = "Total number of compliant devices"
}
