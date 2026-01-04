# Example 1: Remote lock a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Remote lock multiple devices
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_batch" {
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

# Example 3: Remote lock with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_with_validation" {
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

# Example 4: Lock all devices for a specific user (security incident)
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'compromised.user@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_compromised_user" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Lock non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_non_compliant_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 6: Lock iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 7: Emergency lock all corporate Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_corporate_windows" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "locked_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_remote_lock.lock_batch.config.device_ids)
  description = "Number of devices that received remote lock command"
}

output "emergency_locked_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_remote_lock.lock_compromised_user.config.device_ids)
  description = "Number of devices locked in emergency scenario"
}
