# Example 1: Pause configuration refresh on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_single" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 60
      }
    ]
  }
}

# Example 2: Pause configuration refresh on multiple devices with different durations
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_multiple" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 60
      },
      {
        device_id                    = "87654321-4321-4321-4321-ba9876543210"
        pause_time_period_in_minutes = 120
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Pause with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_maximal" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 180
      }
    ]

    comanaged_devices = [
      {
        device_id                    = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        pause_time_period_in_minutes = 90
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Pause configuration refresh during maintenance window
data "microsoft365_graph_beta_device_management_managed_device" "maintenance_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Maintenance Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_maintenance" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.maintenance_devices.items : {
        device_id                    = device.id
        pause_time_period_in_minutes = 240
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
