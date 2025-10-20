---
page_title: "microsoft365_graph_beta_device_management_managed_device_disable Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Disables managed devices from Intune management using the /deviceManagement/managedDevices/{managedDeviceId}/disable and /deviceManagement/comanagedDevices/{managedDeviceId}/disable endpoints. This action disables a device's ability to interact with Intune services while maintaining its enrollment record. Disabled devices cannot receive policies, sync with Intune, or perform managed operations until re-enabled. This is useful for temporarily suspending device management without fully removing the device from Intune, such as during investigations, compliance violations, or security incidents.
  Important Notes:
  Device remains enrolled but cannot sync or receive policiesManagement operations are suspendedDevice can be re-enabled laterLess permanent than retire or wipeUseful for temporary suspensionsSecurity and compliance enforcement
  Use Cases:
  Security incident response (suspected compromise)Compliance violations requiring device suspensionTemporary device quarantineInvestigation of device issuesPreventing policy application during troubleshootingTemporary management suspension
  Platform Support:
  All Platforms: Windows, macOS, iOS/iPadOS, Android
  Reference: Microsoft Graph API - Disable https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_disable (Action)

Disables managed devices from Intune management using the `/deviceManagement/managedDevices/{managedDeviceId}/disable` and `/deviceManagement/comanagedDevices/{managedDeviceId}/disable` endpoints. This action disables a device's ability to interact with Intune services while maintaining its enrollment record. Disabled devices cannot receive policies, sync with Intune, or perform managed operations until re-enabled. This is useful for temporarily suspending device management without fully removing the device from Intune, such as during investigations, compliance violations, or security incidents.

**Important Notes:**
- Device remains enrolled but cannot sync or receive policies
- Management operations are suspended
- Device can be re-enabled later
- Less permanent than retire or wipe
- Useful for temporary suspensions
- Security and compliance enforcement

**Use Cases:**
- Security incident response (suspected compromise)
- Compliance violations requiring device suspension
- Temporary device quarantine
- Investigation of device issues
- Preventing policy application during troubleshooting
- Temporary management suspension

**Platform Support:**
- **All Platforms**: Windows, macOS, iOS/iPadOS, Android

**Reference:** [Microsoft Graph API - Disable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [disable action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta)
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
| **Windows** | ✅ Full Support | Enrolled in Intune |
| **macOS** | ✅ Full Support | Enrolled in Intune |
| **iOS/iPadOS** | ✅ Full Support | Enrolled in Intune |
| **Android** | ✅ Full Support | Enrolled in Intune |

### What is Device Disable?

Device Disable is an action that:
- Temporarily suspends a device's ability to interact with Intune
- Prevents the device from syncing with Intune services
- Blocks policy application and updates
- Maintains the device's enrollment record in Intune
- Preserves all user data on the device
- Can be reversed by re-enabling the device
- Less permanent than retire or wipe actions

### Disable vs Other Management Actions

| Action | Sync Blocked | Policy Blocked | Enrollment Maintained | Data Preserved | Reversible |
|--------|--------------|----------------|----------------------|----------------|------------|
| **Disable** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes (re-enable) |
| **Deprovision** | ❌ No | ✅ Yes | ✅ Yes | ✅ Yes | ⚠️ Partial (re-enroll) |
| **Retire** | ✅ Yes | ✅ Yes | ❌ No | ✅ Yes | ❌ No (new enrollment) |
| **Wipe** | ✅ Yes | ✅ Yes | ❌ No | ❌ No | ❌ No (factory reset) |

### When to Disable Devices

- **Security Incidents**: Suspected device compromise requiring immediate isolation
- **Compliance Violations**: Devices failing to meet compliance requirements
- **Temporary Quarantine**: Investigation of device or user issues
- **Policy Troubleshooting**: Preventing policy application during investigation
- **Management Suspension**: Temporary suspension of management operations
- **Compliance Enforcement**: Enforcing security policies through device suspension

### What Happens When Device is Disabled

- Device is marked as disabled in Intune
- Device cannot sync with Intune services
- New policies and updates are not applied
- Existing policies on device remain in effect
- Device enrollment record is maintained
- User can continue to use the device locally
- All user data and applications remain intact
- Device can be re-enabled to restore management

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to disable. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to disable. These are devices fully managed by Intune only.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to disable different types of devices in one action.

**Important:** Disabled devices will not be able to sync with Intune or receive policy updates until they are re-enabled.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

