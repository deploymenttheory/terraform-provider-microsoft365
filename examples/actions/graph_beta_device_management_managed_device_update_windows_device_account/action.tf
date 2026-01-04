# Example 1: Update Windows device account on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_single" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}

# Example 2: Update multiple Windows device accounts
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_multiple" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      },
      {
        device_id                 = "87654321-4321-4321-4321-ba9876543210"
        device_account_email      = "conference-room-02@company.com"
        password                  = "SecurePassword456!"
        password_rotation_enabled = false
        calendar_sync_enabled     = false
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Update with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_maximal" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]

    comanaged_devices = [
      {
        device_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        device_account_email      = "meeting-room-03@company.com"
        password                  = "SecurePassword789!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Update Surface Hub devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "surface_hubs" {
  filter_type  = "odata"
  odata_filter = "model eq 'Surface Hub'"
}

action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_surface_hubs" {
  config {
    managed_devices = [
      for idx, device in data.microsoft365_graph_beta_device_management_managed_device.surface_hubs.items : {
        device_id                 = device.id
        device_account_email      = format("hub-%02d@company.com", idx + 1)
        password                  = format("SecurePass%03d!", idx + 1)
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
