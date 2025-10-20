# Example 1: Retire a single managed device
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Retire multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Retire devices from a data source query
# First, query for devices that meet certain criteria
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

# Then retire those devices
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_non_compliant_devices" {

  # Extract device IDs from the data source
  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Retire devices with specific operating system
data "microsoft365_graph_beta_device_management_managed_device" "old_ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS' and osVersion startsWith '14'"
}

action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_old_ios_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_ios_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Output examples
output "retired_device_count" {
  value       = length(action.retire_batch.device_ids)
  description = "Number of devices retired in batch operation"
}

output "non_compliant_devices_to_retire" {
  value       = length(action.retire_non_compliant_devices.device_ids)
  description = "Number of non-compliant devices being retired"
}

