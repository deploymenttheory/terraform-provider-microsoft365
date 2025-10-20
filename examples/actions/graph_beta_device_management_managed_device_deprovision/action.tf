# Example 1: Deprovision a single device
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_single" {

  managed_devices {
    device_id           = "12345678-1234-1234-1234-123456789abc"
    deprovision_reason  = "Device being transitioned to new management solution"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Deprovision multiple devices
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_multiple" {

  managed_devices {
    device_id          = "12345678-1234-1234-1234-123456789abc"
    deprovision_reason = "Device repurposing for different department"
  }

  managed_devices {
    device_id          = "87654321-4321-4321-4321-ba9876543210"
    deprovision_reason = "Troubleshooting management issues"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Deprovision devices by user
variable "departing_user_devices" {
  description = "Device IDs for departing user"
  type = map(string)
  default = {
    "device1" = "11111111-1111-1111-1111-111111111111"
    "device2" = "22222222-2222-2222-2222-222222222222"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "user_departure" {

  dynamic "managed_devices" {
    for_each = var.departing_user_devices
    content {
      device_id          = managed_devices.value
      deprovision_reason = "User departure - removing management policies"
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Transition from Intune-only to co-management
data "microsoft365_graph_beta_device_management_managed_device" "transition_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Co-Management Transition'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "comanagement_transition" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.transition_devices.items
    content {
      device_id          = managed_devices.value.id
      deprovision_reason = "Transitioning to co-management with Configuration Manager"
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Deprovision co-managed device
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_comanaged" {

  comanaged_devices {
    device_id          = "abcdef12-3456-7890-abcd-ef1234567890"
    deprovision_reason = "Changing management authority to Configuration Manager only"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 6: Bulk deprovision for management troubleshooting
locals {
  problematic_devices = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222",
    "cccc3333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "troubleshooting" {

  dynamic "managed_devices" {
    for_each = local.problematic_devices
    content {
      device_id          = managed_devices.value
      deprovision_reason = "Management troubleshooting - preparing for re-enrollment"
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 7: Prepare devices for repurposing
data "microsoft365_graph_beta_device_management_managed_device" "repurpose_candidates" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Repurpose Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "repurpose_prep" {

  dynamic "managed_devices" {
    for_each = { for device in data.microsoft365_graph_beta_device_management_managed_device.repurpose_candidates.items : device.id => device }
    content {
      device_id          = managed_devices.key
      deprovision_reason = format("Repurposing device %s for new deployment", managed_devices.value.device_name)
    }
  }

  timeouts = {
    invoke = "30m"
  }
}

# Output examples
output "deprovision_summary" {
  value = {
    managed_count   = length(action.deprovision_multiple.managed_devices)
    comanaged_count = length(action.deprovision_comanaged.comanaged_devices)
  }
  description = "Count of devices deprovisioned"
}

# Important Notes:
# Device Deprovision Features:
# - Removes management policies and profiles
# - Device remains enrolled in Intune
# - User data is preserved
# - Less destructive than wipe or retire
# - Requires reason for auditing
# - Primarily for Windows devices
#
# What is Deprovisioning:
# - Management capability removal action
# - Different from retire or wipe
# - Maintains enrollment status
# - Removes active management
# - Policies and profiles removed
# - Device configuration cleaned up
# - User content preserved
#
# When to Deprovision:
# - Transitioning management solutions
# - Moving to co-management
# - Troubleshooting management issues
# - Preparing for device repurposing
# - Removing management overhead
# - Testing enrollment scenarios
# - Management authority changes
#
# What Happens During Deprovision:
# - Command sent to device
# - Management profiles removed
# - Policies unenrolled
# - Configuration cleaned up
# - Enrollment record maintained
# - User data untouched
# - Device remains registered
#
# Deprovision vs Other Actions:
# - Deprovision: Removes management, keeps enrollment, preserves data
# - Retire: Removes management and enrollment, preserves data
# - Wipe: Removes everything, factory resets device
# - Each serves different purposes
#
# Platform Support:
# - Windows: Primary platform
# - Other platforms: Limited or no support
# - Check platform compatibility
# - Verify action availability
#
# Deprovision Reasons:
# - Required for all operations
# - Used for auditing
# - Track management changes
# - Document decisions
# - Compliance reporting
# - Change management records
#
# Best Practices:
# - Always provide clear reason
# - Document management transitions
# - Plan for re-enrollment if needed
# - Test on pilot devices first
# - Communicate with device users
# - Monitor deprovision success
# - Keep audit trail
#
# Post-Deprovision State:
# - Device still appears in Intune
# - Enrollment status maintained
# - No active management
# - Policies not applied
# - Can be re-managed
# - User can continue using
# - Data intact
#
# Re-enrollment After Deprovision:
# - Device can be re-enrolled
# - Fresh management start
# - New policies applied
# - Clean configuration
# - Previous settings cleared
#
# Use Cases by Scenario:
# - Management solution transition
# - Co-management setup
# - Troubleshooting enrollment
# - Device repurposing
# - Management authority change
# - Policy testing
# - Clean slate for new config
#
# Auditing and Compliance:
# - Deprovision reasons logged
# - Track all management changes
# - Compliance reporting
# - Change management
# - Security audits
# - Document transitions
#
# Troubleshooting:
# - Verify device exists
# - Check enrollment status
# - Ensure device is online
# - Verify permissions
# - Review Intune logs
# - Check for errors
# - Monitor completion
#
# Common Scenarios:
# - User departure cleanup
# - Department transfers
# - Management testing
# - Policy troubleshooting
# - Fresh enrollment prep
# - Management authority shifts
#
# Limitations:
# - Primarily Windows devices
# - Requires enrollment
# - Cannot undo easily
# - May require re-enrollment
# - Check platform support
#
# Related Actions:
# - retire: Full device removal
# - wipe: Factory reset
# - sync_device: Force sync
# - Enrollment actions
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta

