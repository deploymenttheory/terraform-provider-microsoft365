---
page_title: "microsoft365_graph_beta_device_management_managed_device_reenable Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Re-enables previously disabled managed devices in Intune using the /deviceManagement/managedDevices/{managedDeviceId}/reenable and /deviceManagement/comanagedDevices/{managedDeviceId}/reenable endpoints. This action restores a disabled device's ability to interact with Intune services, allowing it to sync and receive policy updates again. Re-enabling is the counterpart to the disable action and restores full management capabilities to devices that were temporarily suspended. This is useful after resolving security incidents, compliance violations, or completing investigations that required temporary device suspension.
  Important Notes:
  Only works on previously disabled devicesRestores sync capability with IntuneRe-enables policy applicationMaintains existing enrollmentReverses the disable actionAll platforms supported
  Use Cases:
  Restoring devices after security investigation completionRe-enabling compliant devices after violations resolvedEnding temporary quarantine periodResuming management after troubleshootingRestoring devices after policy fixesCompleting incident response procedures
  Platform Support:
  All Platforms: Windows, macOS, iOS/iPadOS, Android
  Reference: Microsoft Graph API - Reenable https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_reenable (Action)

Re-enables previously disabled managed devices in Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/reenable` and `/deviceManagement/comanagedDevices/{managedDeviceId}/reenable` endpoints. This action restores a disabled device's ability to interact with Intune services, allowing it to sync and receive policy updates again. Re-enabling is the counterpart to the disable action and restores full management capabilities to devices that were temporarily suspended. This is useful after resolving security incidents, compliance violations, or completing investigations that required temporary device suspension.

**Important Notes:**
- Only works on previously disabled devices
- Restores sync capability with Intune
- Re-enables policy application
- Maintains existing enrollment
- Reverses the disable action
- All platforms supported

**Use Cases:**
- Restoring devices after security investigation completion
- Re-enabling compliant devices after violations resolved
- Ending temporary quarantine period
- Resuming management after troubleshooting
- Restoring devices after policy fixes
- Completing incident response procedures

**Platform Support:**
- **All Platforms**: Windows, macOS, iOS/iPadOS, Android

**Reference:** [Microsoft Graph API - Reenable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [reenable action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Device Management Guides
- [Remote actions in Microsoft Intune](https://learn.microsoft.com/en-us/mem/intune/remote-actions/device-management)
- [Device compliance in Intune](https://learn.microsoft.com/en-us/mem/intune/protect/device-compliance-get-started)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`
- **Delegated**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | Previously disabled device enrolled in Intune |
| **macOS** | ✅ Full Support | Previously disabled device enrolled in Intune |
| **iOS/iPadOS** | ✅ Full Support | Previously disabled device enrolled in Intune |
| **Android** | ✅ Full Support | Previously disabled device enrolled in Intune |

### What is Device Re-enable?

Device Re-enable is an action that:
- Restores a disabled device's ability to interact with Intune
- Re-enables device syncing with Intune services
- Allows policy application and updates to resume
- Maintains the device's enrollment record
- Preserves all user data on the device
- Reverses the effects of the disable action
- Completes the disable-investigate-re-enable workflow

### Disable/Re-enable Workflow

| Step | Action | Device State | Management Active |
|------|--------|--------------|-------------------|
| 1 | Device operating normally | Active | ✅ Yes |
| 2 | Disable action triggered | Disabled | ❌ No |
| 3 | Investigation/resolution occurs | Disabled | ❌ No |
| 4 | **Re-enable action triggered** | **Active** | **✅ Yes** |

### When to Re-enable Devices

- Security investigation has been completed successfully
- Device compliance violations have been resolved
- Temporary quarantine period has ended
- Policy issues have been fixed and tested
- Troubleshooting has been completed
- Incident response procedures are finished
- Device has been cleared to resume normal operations

### What Happens When Device is Re-enabled

- Device is marked as enabled in Intune
- Device can sync with Intune services again
- New policies and updates are applied
- Existing management operations resume
- Device enrollment record remains active
- User can continue using the device normally
- All user data and applications remain intact
- Full management capabilities are restored

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to re-enable. These are devices managed by both Intune and Configuration Manager (SCCM) that were previously disabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to re-enable. These are devices fully managed by Intune that were previously disabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to re-enable different types of devices in one action.

**Important:** Re-enabled devices will be able to sync with Intune and receive policy updates again.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

