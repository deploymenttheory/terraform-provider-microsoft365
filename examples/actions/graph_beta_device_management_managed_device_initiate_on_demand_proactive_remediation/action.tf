# Example 1: Initiate on-demand proactive remediation on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_single" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]
  }
}

# Example 2: Initiate proactive remediation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_multiple" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      },
      {
        device_id        = "abcdef12-3456-7890-abcd-ef1234567890"
        script_policy_id = "11111111-2222-3333-4444-555555555555"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Initiate with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_maximal" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    comanaged_devices = [
      {
        device_id        = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        script_policy_id = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Initiate remediation on all Windows devices with specific script
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_all_windows" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : {
        device_id        = device.id
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "30m"
    }
  }
}
