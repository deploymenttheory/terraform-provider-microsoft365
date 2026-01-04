# Example 1: Initiate MDM key recovery on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_single" {
  config {
    managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Initiate MDM key recovery on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_multiple" {
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

# Example 3: Initiate with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_maximal" {
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

# Example 4: Initiate key recovery on all iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_all_ios" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}
