# Example 1: Rotate local admin password on a single Windows device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Rotate local admin passwords on multiple Windows devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_multiple" {
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

# Example 3: Rotate local admin passwords with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_with_validation" {
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

# Example 4: Rotate local admin passwords on all Windows 10/11 devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 5: Rotate local admin passwords for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_comanaged" {
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

# Example 6: Scheduled rotation for compliance
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "scheduled_rotation" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

    ignore_partial_failures = false

    timeouts = {
      invoke = "25m"
    }
  }
}

# Output examples
output "rotated_passwords_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password.rotate_multiple.config.managed_device_ids)
  description = "Number of devices that had local admin passwords rotated"
}
