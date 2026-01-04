# Example 1: Remove DFCI management from a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Remove DFCI management from multiple devices
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_multiple" {
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

# Example 3: Remove with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_maximal" {
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

# Example 4: Remove DFCI from all Surface devices
data "microsoft365_graph_beta_device_management_managed_device" "surface_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(model, 'Surface')"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_all_surface" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.surface_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 5: Remove DFCI from devices being decommissioned
data "microsoft365_graph_beta_device_management_managed_device" "decommission_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Decommission Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_decommission" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.decommission_devices.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}
