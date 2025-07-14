# Windows Autopilot Device CSV Import and Registration Example
#
# This example demonstrates how to:
# 1. Import Windows Autopilot devices from a CSV file using an ephemeral resource
# 2. Register those devices as Windows Autopilot device identities using for_each

# Configure the Microsoft365 Provider
provider "microsoft365" {
  # Configuration options are omitted for brevity
  # Refer to the provider documentation for configuration details
}

# Import Windows Autopilot devices from a CSV file
ephemeral "microsoft365_graph_beta_windows_autopilot_device_csv_import" "example" {
  file_path = "${path.module}/sample_devices.csv"
}

# Store the imported devices in a local variable for easier reference
locals {
  # Get all devices from the ephemeral resource
  imported_devices = ephemeral.microsoft365_graph_beta_windows_autopilot_device_csv_import.example.devices

  # Create a map with serial number as key for use with for_each
  devices_map = {
    for device in local.imported_devices :
    device.serial_number => device
  }

  # Create a map of devices that should have user assignments
  devices_with_users = {
    for serial, device in local.devices_map :
    serial => device
    if device.assigned_user != ""
  }

  # Create a map of devices without user assignments
  devices_without_users = {
    for serial, device in local.devices_map :
    serial => device
    if device.assigned_user == ""
  }
}

# Register all devices with user assignments using for_each
resource "microsoft365_graph_beta_device_management_windows_autopilot_device_identity" "with_users" {
  for_each = local.devices_with_users

  # Required field - must be unique per device
  serial_number = each.value.serial_number

  # Optional fields for device identification and organization
  group_tag    = each.value.group_tag
  product_key  = each.value.windows_product_id
  display_name = "Autopilot-${each.value.serial_number}"

  # Hardware hash is required for Autopilot registration
  hardware_hash = each.value.hardware_hash

  # User assignment configuration - enables personalized setup
  user_assignment {
    user_principal_name = each.value.assigned_user
  }

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }

  # This resource depends on the ephemeral resource to ensure the devices are imported first
  depends_on = [
    ephemeral.microsoft365_graph_beta_windows_autopilot_device_csv_import.example
  ]
}

# Register all devices without user assignments using for_each
resource "microsoft365_graph_beta_device_management_windows_autopilot_device_identity" "without_users" {
  for_each = local.devices_without_users

  # Required field - must be unique per device
  serial_number = each.value.serial_number

  # Optional fields for device identification and organization
  group_tag    = each.value.group_tag
  product_key  = each.value.windows_product_id
  display_name = "Shared-${each.value.serial_number}"

  # Hardware hash is required for Autopilot registration
  hardware_hash = each.value.hardware_hash

  # No user_assignment block means no user will be assigned to this device

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }

  # This resource depends on the ephemeral resource to ensure the devices are imported first
  depends_on = [
    ephemeral.microsoft365_graph_beta_windows_autopilot_device_csv_import.example
  ]
}

# Output the total number of devices imported and registered
output "total_devices_imported" {
  value       = length(local.imported_devices)
  description = "Total number of devices imported from the CSV file"
}

output "devices_with_users_count" {
  value       = length(local.devices_with_users)
  description = "Number of devices registered with user assignments"
}

output "devices_without_users_count" {
  value       = length(local.devices_without_users)
  description = "Number of devices registered without user assignments"
}

# Output the serial numbers of registered devices
output "registered_devices_with_users" {
  value       = [for device in microsoft365_graph_beta_device_management_windows_autopilot_device_identity.with_users : device.serial_number]
  description = "Serial numbers of registered devices with user assignments"
}

output "registered_devices_without_users" {
  value       = [for device in microsoft365_graph_beta_device_management_windows_autopilot_device_identity.without_users : device.serial_number]
  description = "Serial numbers of registered devices without user assignments"
}

# Output the group tags and their device counts
output "devices_by_group_tag" {
  value = {
    for group_tag in distinct([for device in local.imported_devices : device.group_tag if device.group_tag != ""]) :
    group_tag => length([for device in local.imported_devices : device if device.group_tag == group_tag])
  }
  description = "Count of devices by group tag"
} 