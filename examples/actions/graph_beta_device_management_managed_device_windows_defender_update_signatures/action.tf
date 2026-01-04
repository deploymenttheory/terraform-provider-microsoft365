# Example 1: Update signatures on a single managed device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Update signatures on multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_multiple_managed" {
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

# Example 3: Update signatures on co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_comanaged" {
  config {
    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "11111111-2222-3333-4444-555555555555"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 4: Update both managed and co-managed devices - Maximal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_mixed_devices" {
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
      invoke = "15m"
    }
  }
}

# Example 5: Update all Windows devices using datasource
data "microsoft365_graph_beta_device_management_managed_device" "all_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Update signatures before scheduled scan
data "microsoft365_graph_beta_device_management_managed_device" "workstations" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'WKSTN-')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "pre_scan_update" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 7: Update non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_non_compliant" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 8: Emergency threat response across fleet
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managementAgent eq 'configurationManagerClientMdm')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "emergency_threat_response" {
  config {
    managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : device.id]
    comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "60m"
    }
  }
}
