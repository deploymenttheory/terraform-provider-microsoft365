# Example 1: Set device name for a single device
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_name_single" {

  managed_devices {
    device_id   = "12345678-1234-1234-1234-123456789abc"
    device_name = "NYC-Marketing-Laptop-01"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Set device names for multiple devices with naming convention
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_names_multiple" {

  managed_devices {
    device_id   = "12345678-1234-1234-1234-123456789abc"
    device_name = "NYC-Floor3-Conf-Room-01"
  }

  managed_devices {
    device_id   = "87654321-4321-4321-4321-ba9876543210"
    device_name = "NYC-Floor3-Conf-Room-02"
  }

  managed_devices {
    device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
    device_name = "NYC-Floor3-Conf-Room-03"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Rename devices based on user assignment
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'john.doe@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_user_devices" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.user_devices.items
    content {
      device_id   = managed_devices.value.id
      device_name = format("JohnDoe-%s-%s", managed_devices.value.operating_system, managed_devices.value.serial_number)
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Standardize naming for devices by department
data "microsoft365_graph_beta_device_management_managed_device" "it_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'IT Department'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "standardize_it_devices" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.it_devices.items
    content {
      device_id   = managed_devices.value.id
      device_name = format("IT-DEPT-%s", substr(managed_devices.value.id, 0, 8))
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Rename devices by location
locals {
  device_locations = {
    "12345678-1234-1234-1234-123456789abc" = "NYC-Office"
    "87654321-4321-4321-4321-ba9876543210" = "LA-Office"
    "abcdef12-3456-7890-abcd-ef1234567890" = "Chicago-Office"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_by_location" {

  dynamic "managed_devices" {
    for_each = local.device_locations
    content {
      device_id   = managed_devices.key
      device_name = format("%s-Device-%s", managed_devices.value, formatdate("YYYYMMDDhhmmss", timestamp()))
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Set name for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_comanaged_name" {

  comanaged_devices {
    device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
    device_name = "SCCM-Intune-Hybrid-01"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Rename devices after asset reassignment
variable "reassigned_devices" {
  description = "Map of device IDs to new names after reassignment"
  type = map(string)
  default = {
    "11111111-1111-1111-1111-111111111111" = "Finance-Laptop-A"
    "22222222-2222-2222-2222-222222222222" = "Finance-Laptop-B"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "reassign_devices" {

  dynamic "managed_devices" {
    for_each = var.reassigned_devices
    content {
      device_id   = managed_devices.key
      device_name = managed_devices.value
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Set device names based on model and user
data "microsoft365_graph_beta_device_management_managed_device" "windows_laptops" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceType eq 'desktop'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_windows_laptops" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.windows_laptops.items
    content {
      device_id   = managed_devices.value.id
      device_name = format("WIN-%s-%s", managed_devices.value.model, managed_devices.value.user_display_name)
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "devices_renamed_count" {
  value       = length(action.set_names_multiple.managed_devices)
  description = "Number of devices that had their names set"
}

output "device_naming_info" {
  value = {
    managed   = length(action.set_names_multiple.managed_devices)
    comanaged = length(action.set_comanaged_name.comanaged_devices)
  }
  description = "Count of renamed devices by type"
}

# Important Notes:
# Device Naming Features:
# - Supports all Intune-managed device platforms
# - Each platform may have different naming requirements
# - Device name changes apply after device check-in
# - Custom names improve device identification in console
# - Useful for implementing naming conventions
# - Can rename individual or bulk devices
#
# When to Set Device Names:
# - Implementing organizational naming standards
# - After device reassignment to new users
# - Relocating devices to different offices
# - Standardizing existing device names
# - Organizing devices by department or function
# - Making devices easier to identify and manage
#
# What Happens When Name is Set:
# - Device receives rename command from Intune
# - Name change applies after next device check-in
# - Change reflects in Intune admin console
# - May take minutes to hours depending on check-in
# - Device must be online to receive command
# - User may need to restart device on some platforms
#
# Platform-Specific Considerations:
# - Windows: Computer name changes, may require restart
# - iOS/iPadOS: Device name changes (supervised devices)
# - macOS: Computer name changes for managed devices
# - Android: Varies by enrollment type and management mode
# - Each platform has character and length restrictions
#
# Best Practices:
# - Use consistent naming conventions organization-wide
# - Include identifying information (location, user, function)
# - Avoid special characters that may not be supported
# - Keep names within platform-specific length limits
# - Document your naming convention for consistency
# - Test naming on pilot devices first
# - Consider automation for large-scale renames
#
# Naming Convention Examples:
# - Location-based: "NYC-Floor3-Laptop-01"
# - User-based: "JDoe-MacBook-Pro"
# - Department-based: "IT-Desktop-5"
# - Function-based: "ConferenceRoom-A-iPad"
# - Asset-based: "ASSET-12345"
# - Hybrid: "NYC-IT-JDoe-Laptop"
#
# Character Restrictions (Platform-specific):
# - Windows: 15 characters for NetBIOS, 63 for DNS
# - macOS: Typically no strict limits
# - iOS/iPadOS: Generally flexible naming
# - Avoid special characters: / \ : * ? " < > |
# - Some platforms restrict spaces or require alphanumeric
#
# Security Considerations:
# - Device names may be visible to users
# - Avoid including sensitive information
# - Don't use names that reveal security details
# - Consider privacy when using user names
# - Requires privileged operations permission
# - Audit device name changes for compliance
#
# Related Actions:
# - Device enrollment: Set names during enrollment
# - Bulk operations: Use for mass device renames
# - Inventory management: Track devices by name
# - Compliance: Enforce naming standards
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta

