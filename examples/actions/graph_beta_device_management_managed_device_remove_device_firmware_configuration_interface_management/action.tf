# Example 1: Remove DFCI management from single device
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_dfci_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Remove DFCI from multiple Surface devices
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_dfci_multiple_surface" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Remove DFCI from devices being decommissioned
variable "decommissioned_devices" {
  description = "Device IDs being decommissioned from DFCI management"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "decommission_dfci" {
  managed_device_ids = var.decommissioned_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Remove DFCI based on data source filter
data "microsoft365_graph_beta_device_management_managed_device" "dfci_devices_to_remove" {
  filter_type  = "odata"
  odata_filter = "model eq 'Surface Pro' and deviceCategoryDisplayName eq 'Remove DFCI'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "filtered_removal" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dfci_devices_to_remove.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Transition to standard management
locals {
  transition_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "transition_standard" {
  managed_device_ids = local.transition_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Remove DFCI from co-managed device
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_comanaged_dfci" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Remove DFCI before device transfer
data "microsoft365_graph_beta_device_management_managed_device" "transfer_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Transfer'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "pre_transfer_removal" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.transfer_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Remove DFCI from specific device models
locals {
  surface_models_map = {
    "surface_pro_7" = "11111111-1111-1111-1111-111111111111"
    "surface_pro_8" = "22222222-2222-2222-2222-222222222222"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "surface_models" {
  managed_device_ids = values(local.surface_models_map)

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "dfci_removal_summary" {
  value = {
    managed   = length(action.remove_dfci_multiple_surface.managed_device_ids)
    comanaged = length(action.remove_comanaged_dfci.comanaged_device_ids)
  }
  description = "Count of devices with DFCI removed"
}

