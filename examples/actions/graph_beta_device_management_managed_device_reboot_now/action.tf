# Example 1: Reboot a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Reboot multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_batch" {
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

# Example 3: Reboot with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_with_validation" {
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

# Example 4: Reboot Windows devices with non-compliant state
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_windows_noncompliant" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Reboot kiosk devices
data "microsoft365_graph_beta_device_management_managed_device" "kiosk_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'KIOSK-')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_kiosks" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.kiosk_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Reboot corporate Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_corporate_windows" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 7: Scheduled maintenance reboot for lab devices
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'LAB-')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_lab_maintenance" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "rebooted_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reboot_now.reboot_batch.config.device_ids)
  description = "Number of devices that received reboot command"
}

output "windows_noncompliant_reboot_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reboot_now.reboot_windows_noncompliant.config.device_ids)
  description = "Number of non-compliant Windows devices rebooted"
}
