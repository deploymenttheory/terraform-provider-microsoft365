# Example 1: Disable a single device
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Disable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_multiple" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Disable devices due to security incident
variable "compromised_devices" {
  description = "Device IDs suspected of compromise"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "security_incident" {
  managed_device_ids = var.compromised_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Disable non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "compliance_enforcement" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Temporary quarantine for investigation
locals {
  investigation_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "investigation_quarantine" {
  managed_device_ids = local.investigation_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Disable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_comanaged" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Disable devices by department during policy change
data "microsoft365_graph_beta_device_management_managed_device" "finance_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Finance'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "finance_policy_change" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Disable devices with specific policy violations
data "microsoft365_graph_beta_device_management_managed_device" "policy_violations" {
  filter_type  = "odata"
  odata_filter = "complianceGracePeriodExpirationDateTime lt 2025-01-01T00:00:00Z"
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "policy_enforcement" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.policy_violations.items : device.id]

  timeouts = {
    invoke = "25m"
  }
}

# Output examples
output "disabled_devices_count" {
  value = {
    managed   = length(action.disable_multiple.managed_device_ids)
    comanaged = length(action.disable_comanaged.comanaged_device_ids)
  }
  description = "Count of devices disabled"
}

# Important Notes:
# Device Disable Features:
# - Prevents device from syncing with Intune
# - Blocks policy application
# - Maintains enrollment record
# - Can be re-enabled later
# - Less permanent than retire/wipe
# - All platforms supported
#
# What is Disabling:
# - Temporary management suspension
# - Device remains enrolled
# - Cannot receive policies
# - Cannot sync with Intune
# - Management operations blocked
# - User can still use device
# - Data remains intact
#
# When to Disable Devices:
# - Security incident response
# - Compliance violations
# - Temporary quarantine
# - Investigation of issues
# - Policy application prevention
# - Troubleshooting
# - Temporary suspension needed
#
# What Happens When Disabled:
# - Device cannot sync
# - Policy updates blocked
# - Management operations suspended
# - Enrollment maintained
# - User data preserved
# - Device remains registered
# - Can be re-enabled
#
# Disable vs Other Actions:
# - Disable: Temporary suspension, can re-enable
# - Deprovision: Removes management, keeps enrollment
# - Retire: Full removal from management
# - Wipe: Factory reset device
# - Each serves different purposes
#
# Platform Support:
# - Windows: Fully supported
# - macOS: Fully supported
# - iOS/iPadOS: Fully supported
# - Android: Fully supported
# - All platforms can be disabled
#
# Best Practices:
# - Use for temporary suspensions
# - Document reason for disabling
# - Plan for re-enabling
# - Monitor disabled devices
# - Communicate with users
# - Set timeframes for review
# - Track all disable actions
#
# Re-enabling Devices:
# - Separate action needed
# - Devices can be re-enabled
# - Management resumes
# - Policies re-applied
# - Sync restored
# - Normal operations resume
#
# Security Use Cases:
# - Suspected compromise
# - Policy violations
# - Investigation quarantine
# - Incident response
# - Compliance enforcement
# - Temporary lockout
#
# Compliance Use Cases:
# - Non-compliant devices
# - Grace period expired
# - Policy violations
# - Certificate expiration
# - Security baseline failures
# - Audit requirements
#
# Troubleshooting Use Cases:
# - Policy conflicts
# - Sync issues
# - Configuration problems
# - Testing scenarios
# - Isolation for diagnosis
#
# Post-Disable State:
# - Device shows as disabled
# - Cannot communicate with Intune
# - Policies not enforced
# - User can use device
# - Apps function normally
# - Data intact
# - Enrollment record exists
#
# Auditing and Tracking:
# - Log all disable actions
# - Document reasons
# - Track duration
# - Monitor re-enablement
# - Compliance reporting
# - Security audits
# - Change management
#
# Troubleshooting:
# - Verify device exists
# - Check enrollment status
# - Ensure permissions correct
# - Review Intune logs
# - Verify device ID
# - Check for errors
# - Monitor completion
#
# Common Scenarios:
# - Security incidents
# - Compliance enforcement
# - Policy testing
# - Investigation needs
# - Temporary suspension
# - Department transitions
#
# Limitations:
# - Requires enrollment
# - Cannot sync while disabled
# - No policy updates
# - Management blocked
# - Needs manual re-enable
#
# Related Actions:
# - re-enable: Restore management
# - deprovision: Remove management
# - retire: Full device removal
# - wipe: Factory reset
# - sync_device: Force sync (when enabled)
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta

