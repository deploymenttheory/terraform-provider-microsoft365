---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys Action"
subcategory: "Device Management"
description: |-
  Rotates BitLocker encryption recovery keys on Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/rotateBitLockerKeys and /deviceManagement/comanagedDevices/{managedDeviceId}/rotateBitLockerKeys endpoints. This action generates new BitLocker recovery keys and escrows them to Intune, invalidating the previous recovery keys.
  What This Action Does:
  Generates new BitLocker recovery passwordEscrows new recovery key to Intune/Azure ADInvalidates previous recovery keysUpdates key protector on deviceMaintains encryption state (no re-encryption)Audits key rotation event
  When to Use:
  Security incident or breach responseRecovery key compromised or exposedCompliance policy requirementRegular security maintenance scheduleDevice ownership transferAdministrative access changeProactive security hardening
  Platform Support:
  Windows 10: Pro, Enterprise, Education (v1703+)Windows 11: All editions with BitLockerOther platforms: Not applicable (no BitLocker)
  Important Considerations:
  Only rotates recovery keys, not encryption keysDevice must be online and connectedBitLocker must be enabled and configuredPrevious recovery keys become invalidNew keys escrowed to Azure AD/IntuneNo user interaction requiredNo device restart needed
  Reference: Microsoft Graph API - Rotate BitLocker Keys https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatebitlockerkeys?view=graph-rest-beta
---

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys (Action)

Rotates BitLocker encryption recovery keys on Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/rotateBitLockerKeys` and `/deviceManagement/comanagedDevices/{managedDeviceId}/rotateBitLockerKeys` endpoints. This action generates new BitLocker recovery keys and escrows them to Intune, invalidating the previous recovery keys.

**What This Action Does:**
- Generates new BitLocker recovery password
- Escrows new recovery key to Intune/Azure AD
- Invalidates previous recovery keys
- Updates key protector on device
- Maintains encryption state (no re-encryption)
- Audits key rotation event

**When to Use:**
- Security incident or breach response
- Recovery key compromised or exposed
- Compliance policy requirement
- Regular security maintenance schedule
- Device ownership transfer
- Administrative access change
- Proactive security hardening

**Platform Support:**
- **Windows 10**: Pro, Enterprise, Education (v1703+)
- **Windows 11**: All editions with BitLocker
- **Other platforms**: Not applicable (no BitLocker)

**Important Considerations:**
- Only rotates recovery keys, not encryption keys
- Device must be online and connected
- BitLocker must be enabled and configured
- Previous recovery keys become invalid
- New keys escrowed to Azure AD/Intune
- No user interaction required
- No device restart needed

**Reference:** [Microsoft Graph API - Rotate BitLocker Keys](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatebitlockerkeys?view=graph-rest-beta)

## Use Cases

This action is critical for maintaining BitLocker encryption security and compliance across Windows devices:

## API Documentation

- [Microsoft Graph API - Rotate BitLocker Keys](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatebitlockerkeys?view=graph-rest-beta)

## Permissions

The following Microsoft Graph API permissions are required to use this action:

| Permission Type | Permissions (Least Privileged) |
|:----------------|:------------------------------|
| Delegated (work or school account) | DeviceManagementConfiguration.ReadWrite.All, DeviceManagementManagedDevices.ReadWrite.All |
| Delegated (personal Microsoft account) | Not supported |
| Application | DeviceManagementConfiguration.ReadWrite.All, DeviceManagementManagedDevices.ReadWrite.All |

~> **Note:** This action requires both device configuration and device management write permissions as it modifies BitLocker encryption settings.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Related Documentation

- [Microsoft Intune Remote Actions - Windows](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [Microsoft Intune Remote Actions - iOS/iPadOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [Microsoft Intune Remote Actions - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Microsoft Intune Remote Actions - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [Microsoft Intune Remote Actions - ChromeOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

## Notes

### Platform Compatibility

This remote action is only available for Windows devices with BitLocker encryption enabled. The table below shows platform support:

| Platform | Supported | Notes |
|:---------|:----------|:------|
| **Windows 10** | ✅ | Pro, Enterprise, Education (Version 1703 or later) |
| **Windows 11** | ✅ | All editions with BitLocker support |
| **Windows Server** | ⚠️ | Limited support - depends on BitLocker availability |
| **macOS** | ❌ | Not supported - uses FileVault, not BitLocker |
| **iOS** | ❌ | Not supported - no BitLocker equivalent |
| **iPadOS** | ❌ | Not supported - no BitLocker equivalent |
| **Android** | ❌ | Not supported - uses device encryption, not BitLocker |
| **Android Enterprise** | ❌ | Not supported - uses device encryption, not BitLocker |
| **ChromeOS** | ❌ | Not supported - uses verified boot, not BitLocker |

### Important Considerations

**Key Rotation Fundamentals:**
- **Recovery Keys Only**: This action rotates BitLocker **recovery passwords/keys**, NOT the full-volume encryption keys
- **No Re-Encryption**: Data on the drive is NOT re-encrypted; the operation completes quickly (typically seconds to minutes)
- **Previous Keys Invalid**: All previous recovery keys become immediately invalid and cannot be used for recovery
- **New Key Escrow**: New recovery keys are automatically escrowed to both Intune and Azure AD
- **Cannot Be Undone**: Key rotation is permanent; there is no way to restore previous recovery keys
- **Multiple Drives**: Rotates keys for all BitLocker-protected drives on the device (OS, fixed data, removable drives if configured)

**Device Requirements:**
- **BitLocker Enabled**: Devices must have BitLocker encryption fully enabled and operational
- **Online Status**: Devices must be online and connected to Intune to receive the rotation command
- **Network Connectivity**: Devices need network access to communicate with Intune and Azure AD
- **TPM Chip**: Devices typically require TPM 1.2 or 2.0 chip (depending on Windows version and policy)
- **Windows Edition**: Must be Pro, Enterprise, or Education edition (Home edition does not support BitLocker)
- **Management State**: Devices must be actively enrolled and managed by Intune (not stale or orphaned)

**Operational Impact:**
- **No User Impact**: Key rotation is transparent to end users; no notifications or prompts
- **No Restart Required**: Devices do NOT need to be restarted for key rotation to complete
- **No Performance Impact**: Minimal to no performance impact during rotation (no data re-encryption)
- **Quick Operation**: Rotation typically completes in seconds to a few minutes per device
- **Background Process**: Rotation runs as a background system process without user interaction
- **Session Continuity**: User sessions remain active; no interruption to work

**Key Retrieval & Access:**
- **Intune Admin Center**: New keys viewable at Devices > All devices > [device name] > Recovery keys
- **Azure AD Portal**: New keys viewable at Devices > All devices > [device name] > BitLocker keys
- **Azure AD PowerShell**: Keys retrievable via `Get-AzureADDeviceBitLockerKey` cmdlet
- **Microsoft Graph API**: Keys accessible via `/deviceManagement/managedDevices/{id}/recoveryKeys` endpoint
- **RBAC Permissions**: Requires appropriate permissions to view recovery keys in portals
- **Audit Logging**: Key access is logged in Azure AD and Intune audit logs

**BitLocker Recovery Scenarios:**
- **When Keys Needed**: Recovery keys required when BitLocker enters recovery mode (forgotten password, hardware changes, etc.)
- **Key Format**: Recovery keys are 48-digit numerical passwords (8 groups of 6 digits)
- **Key ID**: Each key has a unique Key ID (GUID) to identify which key protector it belongs to
- **Multiple Protectors**: Devices may have multiple key protectors; rotation affects recovery password protector
- **TPM Protector**: TPM-based protectors (primary unlock method) are NOT affected by this action

**Co-Management Considerations:**
- **Dual Escrow**: Co-managed devices can escrow keys to both Intune and Configuration Manager
- **SCCM Integration**: Ensure Configuration Manager BitLocker management policies align with Intune
- **Authority Conflicts**: Verify which system (Intune or SCCM) has authority for BitLocker management
- **Key Synchronization**: Keys may take time to synchronize between Intune and Configuration Manager
- **Policy Precedence**: Understand which policies take precedence in co-management scenarios

**Security Best Practices:**
- **Regular Rotation Schedule**: Implement quarterly or bi-annual key rotation for all managed Windows devices
- **Document Procedures**: Maintain documented processes for key rotation and emergency key retrieval
- **Test Key Retrieval**: Regularly test the process of retrieving and using recovery keys from Intune/Azure AD
- **Access Controls**: Limit access to BitLocker recovery keys to authorized security and helpdesk personnel only
- **Audit Reviews**: Regularly review audit logs for unauthorized key access attempts
- **Incident Response**: Include key rotation in security incident response playbooks
- **Key Backup**: Ensure key escrow is functioning before performing mass rotations
- **Staged Approach**: For large-scale rotations (500+ devices), use a phased rollout approach
- **Device Validation**: Verify device online status and BitLocker health before initiating rotation
- **Post-Rotation Verification**: Confirm new keys are properly escrowed after rotation completes

**Troubleshooting Common Issues:**
- **Offline Devices**: Rotation fails if device is offline; device will receive command when it reconnects
- **BitLocker Not Enabled**: Rotation fails if BitLocker is not enabled; check BitLocker status first
- **TPM Issues**: Rotation may fail if TPM is locked, disabled, or in maintenance mode
- **Policy Conflicts**: Conflicting BitLocker policies can prevent successful key rotation
- **MBAM Conflicts**: Legacy MBAM (Microsoft BitLocker Administration and Monitoring) configurations may interfere
- **Key Escrow Failures**: Network or connectivity issues can prevent new keys from being escrowed properly
- **Permissions Issues**: Insufficient permissions can cause rotation to fail silently
- **Timing**: Allow several minutes for rotation to complete on busy or slow devices

**Compliance & Audit:**
- **Audit Trail**: All key rotation operations are logged in Intune and Azure AD audit logs
- **Compliance Reporting**: Use Intune compliance reports to track BitLocker encryption status
- **Key Age Tracking**: Monitor when keys were last rotated to ensure compliance with rotation policies
- **Evidence Collection**: Audit logs provide evidence of key rotation for compliance audits
- **Retention**: Ensure audit log retention meets your compliance requirements (typically 90-180 days)

## Example Usage

```terraform
# Data source to find all Windows devices with BitLocker enabled
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter = "operatingSystem eq 'Windows'"
}

# Example 1: Rotate BitLocker keys on specific Windows managed devices
# Use this for targeted key rotation on specific devices (e.g., after security incident)
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_specific_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]
}

# Example 2: Rotate BitLocker keys on co-managed Windows devices
# Use this for devices managed by both Intune and Configuration Manager
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_comanaged_devices" {
  comanaged_device_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

# Example 3: Rotate BitLocker keys on both managed and co-managed devices
# Use this for mixed device management scenarios
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_mixed_devices" {
  managed_device_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  ]

  comanaged_device_ids = [
    "cccccccc-cccc-cccc-cccc-cccccccccccc"
  ]
}

# Example 4: Rotate BitLocker keys on all Windows 10 devices
# Use this for scheduled maintenance or compliance requirement across all Windows 10 devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_10_devices" {
  filter = "operatingSystem eq 'Windows' and contains(osVersion, '10.')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_all_windows_10" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_10_devices.managed_devices : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 5: Rotate BitLocker keys on Windows 11 Enterprise devices
# Use this for targeting specific Windows editions
data "microsoft365_graph_beta_device_management_managed_device" "windows_11_enterprise" {
  filter = "operatingSystem eq 'Windows' and contains(osVersion, '11.') and skuFamily eq 'Windows.Enterprise'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_windows_11_enterprise" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_11_enterprise.managed_devices : device.id]
}

# Example 6: Rotate BitLocker keys on non-compliant Windows devices
# Use this as part of compliance remediation process
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_windows" {
  filter = "complianceState eq 'noncompliant' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_noncompliant" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_windows.managed_devices : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Rotate BitLocker keys on devices by name pattern
# Use this when you need to target specific departments or device groups
data "microsoft365_graph_beta_device_management_managed_device" "finance_dept_devices" {
  filter = "startswith(deviceName, 'FIN-') and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_finance_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_dept_devices.managed_devices : device.id]
}

# Example 8: Rotate BitLocker keys on recently enrolled devices
# Use this to ensure new devices have keys properly escrowed after initial setup
data "microsoft365_graph_beta_device_management_managed_device" "recently_enrolled_windows" {
  filter = "enrolledDateTime gt 2024-01-01T00:00:00Z and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_new_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.recently_enrolled_windows.managed_devices : device.id]
}

# Example 9: Rotate BitLocker keys on Azure AD joined Windows devices
# Use this to target cloud-native Windows devices specifically
data "microsoft365_graph_beta_device_management_managed_device" "azure_ad_joined_windows" {
  filter = "joinType eq 'azureADJoined' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_azure_ad_joined" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.azure_ad_joined_windows.managed_devices : device.id]
}

# Example 10: Rotate BitLocker keys on corporate-owned Windows devices
# Use this to differentiate between corporate and BYOD devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter = "managedDeviceOwnerType eq 'company' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_corporate_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.managed_devices : device.id]
}

# Example 11: Rotate BitLocker keys with custom extended timeout
# Use this for very large-scale operations (500+ devices)
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_with_extended_timeout" {
  managed_device_ids = [
    "device-id-1",
    "device-id-2",
    "device-id-3"
  ]

  timeouts = {
    invoke = "45m"
  }
}

# Example 12: Scheduled BitLocker key rotation for compliance
# Use this for regular security maintenance (e.g., quarterly key rotation)
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_managed" {
  filter = "operatingSystem eq 'Windows' and managementAgent eq 'mdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "quarterly_key_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_managed.managed_devices : device.id]

  timeouts = {
    invoke = "60m"
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to rotate BitLocker keys on. These are Windows devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Note:** Co-managed devices can have BitLocker keys escrowed to both Intune and Configuration Manager. At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `managed_device_ids` (List of String) List of managed device IDs to rotate BitLocker keys on. These are Windows devices fully managed by Intune only. Each ID must be a valid GUID format. BitLocker recovery keys will be rotated on these devices. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to rotate keys on different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

