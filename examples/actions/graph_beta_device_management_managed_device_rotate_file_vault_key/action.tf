# Example 1: Rotate FileVault key on a single macOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Rotate FileVault keys on multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_multiple" {
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

# Example 3: Rotate FileVault keys with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_with_validation" {
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

# Example 4: Rotate FileVault keys on all macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_all_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 5: Rotate FileVault keys for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_comanaged" {
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

# Example 6: Rotate FileVault keys for company-owned macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "company_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_company_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.company_macos.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "rotated_filevault_keys_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key.rotate_multiple.config.managed_device_ids)
  description = "Number of devices that had FileVault keys rotated"
}
