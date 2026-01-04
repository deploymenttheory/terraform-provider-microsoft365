# Example 1: Shutdown a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Shutdown multiple devices
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_batch" {
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

# Example 3: Shutdown with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_with_validation" {
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

# Example 4: Shutdown lab devices for weekend energy conservation
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'LAB-')"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_lab_weekend" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Emergency shutdown for specific device
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "emergency_shutdown" {
  config {
    device_ids = [
      "12345678-abcd-1234-abcd-123456789def"
    ]

    timeouts = {
      invoke = "2m"
    }
  }
}

# Example 6: Shutdown kiosk devices overnight
data "microsoft365_graph_beta_device_management_managed_device" "kiosk_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'KIOSK-')"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_kiosks_overnight" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.kiosk_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "shutdown_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_shutdown.shutdown_batch.config.device_ids)
  description = "Number of devices that received shutdown command"
}

output "lab_shutdown_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_shutdown.shutdown_lab_weekend.config.device_ids)
  description = "Number of lab devices shut down for energy conservation"
}

# Important Note: 
# Shutdown powers devices OFF completely and requires manual power-on to restart.
# Use with caution. Consider using reboot action if devices need to come back online automatically.
