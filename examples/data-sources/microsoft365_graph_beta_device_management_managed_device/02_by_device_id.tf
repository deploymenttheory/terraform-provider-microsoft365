# Example: Get a specific managed device by ID
data "microsoft365_graph_beta_device_management_managed_device" "by_id" {
  device_id = "00000000-0000-0000-0000-000000000001"
}

# Output: Device information by ID
output "device_by_id_info" {
  value = length(data.microsoft365_graph_beta_device_management_managed_device.by_id.items) > 0 ? {
    id            = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].id
    name          = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].device_name
    os            = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].operating_system
    os_version    = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].os_version
    enrolled      = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].enrolled_date_time
    compliance    = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].compliance_state
    serial_number = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].serial_number
    manufacturer  = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].manufacturer
    model         = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].model
  } : null
  description = "Detailed device information for the specified device ID"
}
