# Example 1: Rotate BitLocker keys on a single Windows device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Rotate BitLocker keys on multiple Windows devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_multiple" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Rotate BitLocker keys with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Rotate BitLocker keys on all Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 5: Rotate BitLocker keys for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_comanaged" {
  config {
    comanaged_device_ids = [
      "11111111-1111-1111-1111-111111111111",
      "22222222-2222-2222-2222-222222222222"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Rotate BitLocker keys for non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_noncompliant" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_windows.items : device.id]

    ignore_partial_failures = false

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "rotated_bitlocker_keys_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys.rotate_multiple.config.managed_device_ids)
  description = "Number of devices that had BitLocker keys rotated"
}
