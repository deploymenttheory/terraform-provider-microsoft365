# Example 1: Re-enable a single device
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Re-enable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_multiple" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Re-enable devices after security investigation
variable "investigated_devices" {
  description = "Device IDs cleared from security investigation"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "post_investigation" {
  managed_device_ids = var.investigated_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Re-enable compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "now_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'compliant' and isEnabled eq false"
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "compliance_restored" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.now_compliant.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Re-enable after incident resolution
locals {
  incident_resolved_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "incident_resolved" {
  managed_device_ids = local.incident_resolved_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Re-enable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_comanaged" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Re-enable devices after policy fix
data "microsoft365_graph_beta_device_management_managed_device" "after_policy_fix" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Policy Fix Complete'"
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "policy_fix_complete" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.after_policy_fix.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Re-enable after quarantine period
locals {
  quarantine_period_devices = {
    "device1" = "11111111-1111-1111-1111-111111111111"
    "device2" = "22222222-2222-2222-2222-222222222222"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "quarantine_ended" {
  managed_device_ids = values(local.quarantine_period_devices)

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "reenabled_devices_count" {
  value = {
    managed   = length(action.reenable_multiple.managed_device_ids)
    comanaged = length(action.reenable_comanaged.comanaged_device_ids)
  }
  description = "Count of devices re-enabled"
}

# Important Notes:
# Device Re-enable Features:
# - Restores sync capability
# - Re-enables policy application
# - Maintains enrollment record
# - Counterpart to disable action
# - All platforms supported
# - Reversible action
#
# What is Re-enabling:
# - Restoration of management
# - Reverse of disable action
# - Device can sync again
# - Policies applied again
# - Management operations resume
# - User continues using device
# - Data remains intact
#
# When to Re-enable Devices:
# - Investigation completed
# - Compliance restored
# - Security incident resolved
# - Policy issues fixed
# - Troubleshooting complete
# - Quarantine period ended
# - Temporary suspension over
#
# What Happens When Re-enabled:
# - Device can sync
# - Policy updates resume
# - Management restored
# - Enrollment maintained
# - User data preserved
# - Device fully functional
# - Normal operations resume
#
# Re-enable vs Other Actions:
# - Re-enable: Reverses disable
# - Requires previous disable
# - Restores management
# - Simple and quick
# - Non-destructive
#
# Platform Support:
# - Windows: Fully supported
# - macOS: Fully supported
# - iOS/iPadOS: Fully supported
# - Android: Fully supported
# - All platforms can be re-enabled
#
# Best Practices:
# - Document re-enable reason
# - Verify issues resolved
# - Test device functionality
# - Communicate with users
# - Monitor after re-enable
# - Track all re-enable actions
# - Audit trail maintained
#
# Pre-Requisites:
# - Device must be disabled first
# - Cannot re-enable active device
# - Device must exist in Intune
# - Enrollment must be valid
# - Proper permissions required
#
# Post-Re-enable State:
# - Device shows as enabled
# - Can communicate with Intune
# - Policies enforced
# - User can use normally
# - Apps function normally
# - Data intact
# - Management active
#
# Security Use Cases:
# - Incident resolved
# - Investigation complete
# - Threat eliminated
# - Compliance restored
# - Policy violations fixed
# - Security measures applied
#
# Compliance Use Cases:
# - Device now compliant
# - Grace period granted
# - Violations resolved
# - Certificates renewed
# - Policies applied
# - Audit complete
#
# Troubleshooting Use Cases:
# - Policy conflicts resolved
# - Sync issues fixed
# - Configuration corrected
# - Testing complete
# - Diagnosis finished
#
# Common Workflows:
# - Disable → Investigate → Re-enable
# - Disable → Fix → Re-enable
# - Disable → Wait → Re-enable
# - Disable → Verify → Re-enable
#
# Auditing and Tracking:
# - Log all re-enable actions
# - Document reasons
# - Track duration disabled
# - Monitor post-re-enable
# - Compliance reporting
# - Security audits
# - Change management
#
# Troubleshooting:
# - Verify device exists
# - Check was disabled
# - Ensure permissions correct
# - Review Intune logs
# - Verify device ID
# - Check for errors
# - Monitor completion
#
# Common Scenarios:
# - Security cleared
# - Compliance achieved
# - Policy fixed
# - Investigation done
# - User reinstated
# - Issue resolved
#
# Limitations:
# - Must be previously disabled
# - Requires enrollment
# - Needs manual action
# - Cannot re-enable twice
#
# Related Actions:
# - disable: Suspend management
# - deprovision: Remove management
# - retire: Full device removal
# - wipe: Factory reset
# - sync_device: Force sync
#
# Disable/Re-enable Lifecycle:
# 1. Device active
# 2. Disable action triggered
# 3. Device suspended
# 4. Issue investigated/resolved
# 5. Re-enable action triggered
# 6. Device active again
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta

