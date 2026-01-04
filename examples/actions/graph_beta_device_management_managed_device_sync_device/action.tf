# Example 1: Sync a single managed device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Sync multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_managed_only" {
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

# Example 3: Sync co-managed devices only
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_comanaged_only" {
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

# Example 4: Sync both managed and co-managed devices - Maximal
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_mixed_devices" {
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

# Example 5: Sync all Windows devices using datasource
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Sync non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_non_compliant" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 7: Sync iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_ios_devices" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 8: Emergency policy deployment
data "microsoft365_graph_beta_device_management_managed_device" "all_managed" {
  filter_type = "all"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClientMdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "emergency_sync_all" {
  config {
    managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_managed.items : device.id]
    comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "60m"
    }
  }
}
