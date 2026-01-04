# Example 1: Clean single Windows device (remove user data) - Minimal
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_full_clean" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 2: Clean single Windows device (preserve user data)
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_preserve_data" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = true
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Clean multiple Windows devices with different options per device
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "multiple_devices_mixed_options" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      },
      {
        device_id      = "87654321-4321-4321-4321-ba9876543210"
        keep_user_data = true
      },
      {
        device_id      = "abcdef12-3456-7890-abcd-ef1234567890"
        keep_user_data = false
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 4: Clean co-managed Windows devices
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "comanaged_devices" {
  config {
    comanaged_devices = [
      {
        device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        keep_user_data = false
      },
      {
        device_id      = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
        keep_user_data = true
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Clean Windows devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : {
        device_id      = device.id
        keep_user_data = false
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Maximal configuration with both managed and co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "maximal_config" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      },
      {
        device_id      = "87654321-4321-4321-4321-987654321cba"
        keep_user_data = true
      }
    ]

    comanaged_devices = [
      {
        device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        keep_user_data = false
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}
