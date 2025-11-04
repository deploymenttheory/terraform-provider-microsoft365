---
page_title: "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Rotates the FileVault recovery key for macOS managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/rotateFileVaultKey and /deviceManagement/comanagedDevices/{managedDeviceId}/rotateFileVaultKey endpoints. This action generates a new FileVault recovery key and escrows it with Intune, ensuring that administrators can recover encrypted macOS devices if users forget their passwords or lose access. Regular key rotation is a security best practice that limits the window of exposure if a key is compromised. This action supports rotating keys on multiple devices in a single operation.
  Important Notes:
  Only applicable to macOS devices with FileVault enabledGenerates a new personal recovery keyNew key is escrowed with Intune automaticallyPrevious recovery key becomes invalidDevice must be online to receive rotation commandUser does not need to be logged inNo user interaction required for rotation
  Use Cases:
  Regular security key rotation complianceAfter potential key compromise or exposureWhen changing device ownership or assignmentAs part of security incident responsePeriodic rotation per security policyBefore device reassignment to new users
  Platform Support:
  macOS: Fully supported on devices with FileVault enabledOther Platforms: Not applicable (FileVault is macOS-only)
  Reference: Microsoft Graph API - Rotate FileVault Key https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key (Action)

Rotates the FileVault recovery key for macOS managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/rotateFileVaultKey` and `/deviceManagement/comanagedDevices/{managedDeviceId}/rotateFileVaultKey` endpoints. This action generates a new FileVault recovery key and escrows it with Intune, ensuring that administrators can recover encrypted macOS devices if users forget their passwords or lose access. Regular key rotation is a security best practice that limits the window of exposure if a key is compromised. This action supports rotating keys on multiple devices in a single operation.

**Important Notes:**
- Only applicable to macOS devices with FileVault enabled
- Generates a new personal recovery key
- New key is escrowed with Intune automatically
- Previous recovery key becomes invalid
- Device must be online to receive rotation command
- User does not need to be logged in
- No user interaction required for rotation

**Use Cases:**
- Regular security key rotation compliance
- After potential key compromise or exposure
- When changing device ownership or assignment
- As part of security incident response
- Periodic rotation per security policy
- Before device reassignment to new users

**Platform Support:**
- **macOS**: Fully supported on devices with FileVault enabled
- **Other Platforms**: Not applicable (FileVault is macOS-only)

**Reference:** [Microsoft Graph API - Rotate FileVault Key](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [rotateFileVaultKey action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device rotate FileVault](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-rotate-filevault)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **macOS** | ✅ Full Support | FileVault must be enabled on device |
| **Windows** | ❌ Not Supported | FileVault is macOS-only (use BitLocker rotation instead) |
| **iOS/iPadOS** | ❌ Not Supported | FileVault is macOS-only |
| **Android** | ❌ Not Supported | FileVault is macOS-only |

### What is FileVault Key Rotation?

FileVault Key Rotation is an action that:
- Generates a new FileVault personal recovery key
- Automatically escrows the new key with Intune
- Invalidates the previous recovery key
- Operates without user interaction
- Enhances security through regular key changes
- Maintains continuous disk encryption protection

### When to Rotate FileVault Keys

- Regular compliance-driven rotation (quarterly/annually per security policy)
- After suspected recovery key compromise or exposure
- When reassigning devices to new users or departments
- As part of security incident response procedures
- Before or after employee termination or transfer
- To meet regulatory or audit requirements
- After key has been accessed by administrative staff

### What Happens When FileVault Key is Rotated

- Intune sends rotation command to the macOS device
- Device generates new unique FileVault recovery key
- New key is automatically escrowed with Intune
- Previous recovery key is invalidated immediately
- Process completes without user interaction or awareness
- No device restart or user password change required
- Disk encryption continues without interruption
- New key becomes available in Intune portal for admin access

## Example Usage

```terraform
# Example 1: Rotate FileVault key for a single macOS device
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Rotate FileVault keys for multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_multiple" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Rotate keys for all managed macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_all_macos" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 4: Rotate keys for macOS devices in specific department
data "microsoft365_graph_beta_device_management_managed_device" "finance_macos" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS' and deviceCategoryDisplayName eq 'Finance'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_finance_macos" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_macos.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Scheduled quarterly key rotation
locals {
  # Devices due for quarterly rotation
  macos_devices_for_rotation = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "quarterly_rotation" {
  managed_device_ids = local.macos_devices_for_rotation

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Rotate key for co-managed macOS device
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_comanaged" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Rotate keys after security incident
variable "compromised_device_ids" {
  description = "List of device IDs potentially affected by security incident"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "incident_response" {
  managed_device_ids = var.compromised_device_ids

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Rotate keys for devices being reassigned
data "microsoft365_graph_beta_device_management_managed_device" "reassignment_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'previous.user@example.com' and operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "reassignment_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.reassignment_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 9: Rotate keys for devices with old enrollment dates
data "microsoft365_graph_beta_device_management_managed_device" "old_macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS' and enrolledDateTime lt 2024-01-01T00:00:00Z"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_old_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_macos_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Output examples
output "devices_rotated_count" {
  value       = length(action.rotate_multiple.managed_device_ids)
  description = "Number of devices that had FileVault keys rotated"
}

output "rotation_summary" {
  value = {
    managed   = length(action.rotate_all_macos.managed_device_ids)
    comanaged = length(action.rotate_comanaged.comanaged_device_ids)
  }
  description = "Count of FileVault key rotations by device type"
}

# Important Notes:
# FileVault Key Rotation Features:
# - Generates new FileVault recovery key
# - New key is automatically escrowed with Intune
# - Previous recovery key becomes invalid
# - No user interaction required
# - Device does not need to be logged in
# - Works silently in the background
# - Requires macOS devices with FileVault enabled
#
# When to Rotate FileVault Keys:
# - Regular compliance-driven key rotation (quarterly/annually)
# - After suspected key compromise or exposure
# - When changing device ownership or assignment
# - As part of security incident response
# - Before device reassignment to new users
# - After employee termination or transfer
# - Compliance requirements mandate key rotation
#
# What Happens When Key is Rotated:
# - Device receives rotation command from Intune
# - macOS generates new FileVault recovery key
# - New key is escrowed with Intune automatically
# - Previous recovery key is invalidated
# - Change occurs without user interaction
# - No device restart required
# - User passwords remain unchanged
# - Disk encryption continues uninterrupted
#
# Platform Requirements:
# - macOS: Fully supported (FileVault-enabled devices)
# - Device must have FileVault encryption enabled
# - Device must be enrolled in Intune
# - Device must be online to receive command
# - Other platforms: Not applicable (FileVault is macOS-only)
#
# Best Practices:
# - Implement regular rotation schedule (quarterly)
# - Rotate keys after employee changes
# - Document rotation policy and schedule
# - Test rotation on pilot devices first
# - Monitor rotation success/failure rates
# - Keep audit logs of all rotations
# - Rotate immediately after suspected compromise
# - Combine with other security measures
#
# Security Benefits:
# - Limits exposure window if key compromised
# - Compliance with security policies
# - Part of defense-in-depth strategy
# - Reduces risk of unauthorized access
# - Maintains control over device recovery
# - Supports zero-trust security model
# - Helps meet regulatory requirements
#
# FileVault Overview:
# - Full disk encryption for macOS
# - Encrypts entire startup disk
# - Recovery key allows admin unlock
# - Required for many compliance standards
# - Protects data at rest
# - Transparent to users
# - Minimal performance impact
#
# Recovery Key Management:
# - Keys escrowed with Intune
# - Accessible to admins via Intune portal
# - Used when user forgets password
# - Required for device recovery
# - Rotation creates new unique key
# - Old key no longer works after rotation
# - Keys stored securely in Intune
#
# Compliance Considerations:
# - Many standards require key rotation
# - NIST recommends periodic rotation
# - Document rotation frequency
# - Maintain rotation audit trail
# - Verify rotation completion
# - Report on rotation compliance
# - Integrate with compliance dashboards
#
# Troubleshooting:
# - Verify device is macOS
# - Check FileVault is enabled
# - Ensure device is online
# - Verify Intune connectivity
# - Check device compliance status
# - Review device logs if rotation fails
# - Contact Microsoft support if needed
#
# Related Actions:
# - Device compliance: Enforce FileVault
# - Encryption policies: Configure FileVault
# - Retire device: Remove escrowed keys
# - Remote lock: Use with rotated keys
# - Device wipe: Secure data removal
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices fully managed by Intune only. Each device must have FileVault encryption enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to rotate keys on different types of devices in one action.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

