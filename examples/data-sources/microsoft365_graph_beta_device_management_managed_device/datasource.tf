# Example 1: Get all managed devices
data "microsoft365_graph_beta_device_management_managed_device" "all" {
  filter_type = "all"
}

# Example 2: Get a specific managed device by ID
data "microsoft365_graph_beta_device_management_managed_device" "by_id" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000000"
}

# Example 3: Get managed devices by device name (partial match)
data "microsoft365_graph_beta_device_management_managed_device" "by_device_name" {
  filter_type  = "device_name"
  filter_value = "DESKTOP"
}

# Example 4: Get managed devices by serial number (partial match)
data "microsoft365_graph_beta_device_management_managed_device" "by_serial_number" {
  filter_type  = "serial_number"
  filter_value = "ABC123"
}

# Example 5: Get managed devices by user ID (partial match)
data "microsoft365_graph_beta_device_management_managed_device" "by_user_id" {
  filter_type  = "user_id"
  filter_value = "user@example.com"
}

# Example 6: Get managed devices using OData filter (Windows devices only)
data "microsoft365_graph_beta_device_management_managed_device" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

# Example 7: Advanced OData query with filter, orderby, and select
data "microsoft365_graph_beta_device_management_managed_device" "odata_advanced" {
  filter_type   = "odata"
  odata_filter  = "operatingSystem eq 'Windows'"
  odata_orderby = "deviceName"
  odata_select  = "id,deviceName,operatingSystem,complianceState"
}

# Example 8: Comprehensive OData query with top and orderby
data "microsoft365_graph_beta_device_management_managed_device" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "operatingSystem eq 'Windows'"
  odata_top     = 50
  odata_orderby = "lastSyncDateTime desc"
}

# Example 9: OData with count and filter
data "microsoft365_graph_beta_device_management_managed_device" "odata_with_count" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'compliant'"
  odata_count  = true
}

# Example 10: OData search query
data "microsoft365_graph_beta_device_management_managed_device" "odata_search" {
  filter_type  = "odata"
  odata_search = "\"displayName:LAPTOP\""
  odata_count  = true
}

# Example 11: OData with expand to include related entities
data "microsoft365_graph_beta_device_management_managed_device" "odata_expand" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
  odata_expand = "deviceCategory"
}

# Output examples
output "all_managed_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.all.items)
  description = "Total number of managed devices"
}

output "windows_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.odata_advanced.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      compliance_state = device.compliance_state
    }
  ]
  description = "List of Windows devices with selected fields"
}

output "compliant_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.odata_with_count.items)
  description = "Number of compliant devices"
}

output "device_by_id_info" {
  value = length(data.microsoft365_graph_beta_device_management_managed_device.by_id.items) > 0 ? {
    name       = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].device_name
    os         = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].operating_system
    enrolled   = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].enrolled_date_time
    compliance = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].compliance_state
  } : null
  description = "Device information by ID"
}
