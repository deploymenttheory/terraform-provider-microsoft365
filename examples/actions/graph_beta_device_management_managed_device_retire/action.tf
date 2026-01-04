# Example 1: Retire a single managed device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Retire multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_batch" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Retire with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_with_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Retire devices from a data source query
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_non_compliant_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Retire devices with specific operating system
data "microsoft365_graph_beta_device_management_managed_device" "old_ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (startsWith(osVersion, '14'))"
}

action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_old_ios_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_ios_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "retired_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_retire.retire_batch.config.device_ids)
  description = "Number of devices retired in batch operation"
}

output "non_compliant_devices_to_retire" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_retire.retire_non_compliant_devices.config.device_ids)
  description = "Number of non-compliant devices being retired"
}
