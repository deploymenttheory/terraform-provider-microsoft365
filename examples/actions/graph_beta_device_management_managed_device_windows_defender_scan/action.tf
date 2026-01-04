# Example 1: Quick scan on a single Windows device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_single" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      }
    ]
  }
}

# Example 2: Full scan on a single Windows device
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_single" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = false
      }
    ]
  }
}

# Example 3: Mixed scans on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "mixed_scans" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      },
      {
        device_id  = "87654321-4321-4321-4321-ba9876543210"
        quick_scan = false
      }
    ]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 4: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_maximal" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      },
      {
        device_id  = "87654321-4321-4321-4321-ba9876543210"
        quick_scan = false
      }
    ]

    comanaged_devices = [
      {
        device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        quick_scan = true
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Quick scan all Windows devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_all_windows" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : {
        device_id  = device.id
        quick_scan = true
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Full scan on non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_windows.items : {
        device_id  = device.id
        quick_scan = false
      }
    ]

    ignore_partial_failures = false

    timeouts = {
      invoke = "60m"
    }
  }
}

# Example 7: Scan co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        quick_scan = true
      },
      {
        device_id  = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
        quick_scan = false
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "scanned_devices_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_windows_defender_scan.mixed_scans.config.managed_devices)
  description = "Number of devices that had scans initiated"
}
