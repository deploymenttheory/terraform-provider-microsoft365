# Example 1: Reset passcode for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Reset passcodes for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_batch" {
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

# Example 3: Reset passcodes with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_with_validation" {
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

# Example 4: Reset passcodes for locked devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "locked_ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_locked_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.locked_ios_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Emergency passcode reset for specific user's devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'user@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_user_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Reset passcodes for Android devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_android" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "reset_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reset_passcode.reset_batch.config.device_ids)
  description = "Number of devices for which passcodes were reset"
}

output "locked_devices_reset_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reset_passcode.reset_locked_devices.config.device_ids)
  description = "Number of locked devices for which passcodes were reset"
}
