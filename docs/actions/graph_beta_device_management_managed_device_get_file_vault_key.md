---
page_title: "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves the FileVault recovery key for macOS managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/getFileVaultKey and /deviceManagement/comanagedDevices/{managedDeviceId}/getFileVaultKey endpoints. This action allows administrators to retrieve escrowed FileVault recovery keys for device recovery purposes. The recovery key is displayed in the action output and can be used to unlock an encrypted macOS device when a user has forgotten their password or is otherwise unable to access the device. This is a critical capability for IT support and device recovery scenarios.
  Important Security Notes:
  Recovery keys are highly sensitive credentialsKeys grant full access to encrypted device dataAccess to keys should be audited and restrictedOnly retrieve keys when necessary for device recoveryKeys are displayed in plain text in action outputEnsure proper security controls on Terraform stateConsider security implications before using in automation
  Use Cases:
  Emergency device recovery when user cannot log inUnlocking devices for departing employeesTechnical support scenarios requiring device accessDisaster recovery and business continuityDevice repurposing or reassignment preparation
  Platform Support:
  macOS: Fully supported on devices with FileVault enabled and keys escrowedOther Platforms: Not applicable (FileVault is macOS-only)
  Reference: Microsoft Graph API - Get FileVault Key https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_get_file_vault_key (Action)

Retrieves the FileVault recovery key for macOS managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/getFileVaultKey` and `/deviceManagement/comanagedDevices/{managedDeviceId}/getFileVaultKey` endpoints. This action allows administrators to retrieve escrowed FileVault recovery keys for device recovery purposes. The recovery key is displayed in the action output and can be used to unlock an encrypted macOS device when a user has forgotten their password or is otherwise unable to access the device. This is a critical capability for IT support and device recovery scenarios.

**Important Security Notes:**
- Recovery keys are highly sensitive credentials
- Keys grant full access to encrypted device data
- Access to keys should be audited and restricted
- Only retrieve keys when necessary for device recovery
- Keys are displayed in plain text in action output
- Ensure proper security controls on Terraform state
- Consider security implications before using in automation

**Use Cases:**
- Emergency device recovery when user cannot log in
- Unlocking devices for departing employees
- Technical support scenarios requiring device access
- Disaster recovery and business continuity
- Device repurposing or reassignment preparation

**Platform Support:**
- **macOS**: Fully supported on devices with FileVault enabled and keys escrowed
- **Other Platforms**: Not applicable (FileVault is macOS-only)

**Reference:** [Microsoft Graph API - Get FileVault Key](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta)

## ⚠️ Security Warning

**This action retrieves and displays FileVault recovery keys in plain text.** Recovery keys are highly sensitive credentials that grant full access to encrypted device data. 

- Only use this action when necessary for legitimate device recovery purposes
- Ensure proper security controls are in place for Terraform state files
- Keys will be displayed in action output and may be stored in state
- Access should be logged, audited, and restricted to authorized personnel
- Follow your organization's security policies for handling sensitive credentials

## Microsoft Documentation

### Graph API References
- [getFileVaultKey function](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune FileVault Guides
- [Use disk encryption for macOS with Intune](https://learn.microsoft.com/en-us/mem/intune/protect/encrypt-devices-filevault)
- [FileVault recovery key rotation](https://learn.microsoft.com/en-us/mem/intune/protect/encrypt-devices-filevault#rotate-recovery-keys)

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
| **macOS** | ✅ Full Support | FileVault enabled with key escrowed to Intune |
| **Windows** | ❌ Not Supported | FileVault is macOS-only (use BitLocker recovery for Windows) |
| **iOS/iPadOS** | ❌ Not Supported | FileVault is macOS-only |
| **Android** | ❌ Not Supported | FileVault is macOS-only |

### What is FileVault Key Retrieval?

FileVault Key Retrieval is an action that:
- Retrieves the personal recovery key for FileVault-encrypted macOS devices
- Returns keys that have been escrowed with Intune during encryption setup
- Allows administrators to unlock devices when users cannot access them
- Displays keys in plain text in the action output
- Does not modify the device or recovery key
- Critical capability for device recovery and support scenarios

### When to Retrieve FileVault Keys

- User is locked out and cannot remember their password
- Device needs to be accessed for emergency data recovery
- Departing employee's device needs to be unlocked for data migration
- Device is being repurposed or reassigned to a new user
- Technical support requires access to diagnose hardware/software issues
- Disaster recovery or business continuity scenario
- Legal, compliance, or audit requirement to access device data
- Device recovery after hardware repair or replacement

### What Happens When Key is Retrieved

- Intune returns the escrowed FileVault personal recovery key
- Key is displayed in action output via progress messages
- No changes are made to the device or the key itself
- Device remains encrypted and in its current state
- Retrieved key remains valid until manually rotated
- User's password and device settings remain unchanged
- Administrator can use the key to unlock the device as needed

### How to Use Retrieved Recovery Keys

1. Boot the macOS device to macOS Recovery (hold Command+R during startup)
2. When prompted, select "Unlock with Recovery Key"
3. Enter the retrieved recovery key exactly as displayed
4. Device will unlock and boot normally
5. User can then reset their password if needed
6. Key can be used multiple times until rotated

## Example Usage

```terraform
# SECURITY WARNING: This action retrieves FileVault recovery keys in plain text.
# Recovery keys are highly sensitive credentials. Use appropriate security controls.

# Example 1: Retrieve FileVault key for a single device
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Retrieve keys for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_multiple" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Emergency device recovery scenario
variable "locked_device_ids" {
  description = "Device IDs that need recovery key retrieval"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "emergency_recovery" {
  managed_device_ids = var.locked_device_ids

  timeouts = {
    invoke = "5m"
  }
}

# Example 4: Retrieve key for departing employee's device
data "microsoft365_graph_beta_device_management_managed_device" "departing_employee" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'departing.user@example.com' and operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "departing_employee_recovery" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_employee.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Retrieve key for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_comanaged" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 6: Retrieve keys for devices being reassigned
data "microsoft365_graph_beta_device_management_managed_device" "reassignment_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Pending Reassignment' and operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "reassignment_recovery" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.reassignment_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 7: Support ticket scenario
locals {
  # Device IDs from support ticket requiring recovery
  support_ticket_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "support_ticket_recovery" {
  managed_device_ids = local.support_ticket_devices

  timeouts = {
    invoke = "10m"
  }
}

# Output examples - NOTE: These will contain sensitive recovery keys
# Consider carefully whether outputs are appropriate for your security requirements
output "recovery_operation_summary" {
  value = {
    managed_devices_count   = length(action.retrieve_multiple.managed_device_ids)
    comanaged_devices_count = length(action.retrieve_comanaged.comanaged_device_ids)
  }
  description = "Summary of recovery key retrieval operation (does not contain actual keys)"
}

# Important Notes:
# FileVault Key Retrieval Features:
# - Retrieves personal recovery key for FileVault-encrypted macOS devices
# - Keys are escrowed with Intune during FileVault setup
# - Allows administrators to unlock devices when users cannot
# - Critical for device recovery and support scenarios
# - Keys displayed in plain text in action output
# - Each device has a unique recovery key
#
# Security Considerations:
# - Recovery keys are HIGHLY SENSITIVE credentials
# - Keys grant full access to all encrypted data on device
# - Only retrieve keys when necessary for legitimate recovery
# - Access should be logged and audited
# - Keys displayed in Terraform output and potentially in state
# - Ensure Terraform state has appropriate security controls
# - Consider regulatory and compliance requirements
# - Limit access to personnel with legitimate need
# - Document key retrieval in change management system
#
# When to Retrieve FileVault Keys:
# - User locked out of device and cannot remember password
# - Device needs to be accessed for emergency data recovery
# - Departing employee's device needs to be unlocked
# - Device being repurposed or reassigned to new user
# - Technical support requires access to encrypted device
# - Disaster recovery or business continuity scenario
# - Legal or compliance requirement to access device data
#
# What Happens When Key is Retrieved:
# - Intune returns the escrowed FileVault recovery key
# - Key is displayed in action output (SendProgress messages)
# - Administrator can use key to unlock the device
# - No change occurs on the actual device
# - Device remains encrypted and in current state
# - Key remains valid until next key rotation
# - User's password remains unchanged
#
# Using Retrieved Recovery Keys:
# - Boot macOS device to recovery screen (Command+R at startup)
# - Select "Unlock with Recovery Key" option
# - Enter the retrieved recovery key when prompted
# - Device will unlock and boot normally
# - User can then reset their password if needed
# - Key can be used multiple times until rotated
#
# Platform Requirements:
# - macOS: Fully supported (FileVault-enabled devices)
# - Device must have FileVault encryption enabled
# - Recovery key must be escrowed with Intune
# - Device must be enrolled in Intune
# - Other platforms: Not applicable (FileVault is macOS-only)
#
# Best Practices:
# - Only retrieve keys when absolutely necessary
# - Document business justification for retrieval
# - Log all key retrieval actions
# - Consider key rotation after retrieval
# - Securely communicate keys to authorized personnel
# - Delete keys from temporary storage after use
# - Verify requestor authorization before retrieval
# - Follow organization's key handling procedures
# - Consider time-limited access to retrieved keys
# - Maintain audit trail of key usage
#
# Terraform State Security:
# - Keys may be stored in Terraform state files
# - Ensure state files are encrypted at rest
# - Use remote state with encryption (e.g., S3 with KMS)
# - Restrict access to state files
# - Consider state file retention policies
# - May want to use separate state for sensitive operations
# - Review state file access logs regularly
#
# Alternative Approaches:
# - Use Intune portal for ad-hoc key retrieval
# - Microsoft Endpoint Manager admin center
# - PowerShell with Microsoft Graph
# - Consider if automation is appropriate for your security model
# - Manual retrieval may be more appropriate for high-security environments
#
# Compliance and Auditing:
# - Document who retrieved keys and when
# - Maintain audit log of key access
# - Track business justification for each retrieval
# - Report key retrieval metrics to security team
# - Ensure compliance with data protection regulations
# - May require approval workflow before retrieval
# - Consider integration with SIEM/logging systems
#
# Troubleshooting:
# - Verify device is macOS
# - Check FileVault is enabled on device
# - Ensure key has been escrowed to Intune
# - Verify appropriate API permissions
# - Check device enrollment status
# - Review Intune portal for device details
# - Confirm device is not in error state
#
# Related Actions:
# - rotate_file_vault_key: Generate new recovery key
# - retire: Remove device from management (removes escrowed key)
# - remote_lock: Lock device remotely
# - wipe: Securely erase device data
#
# FileVault Background:
# - Full disk encryption for macOS
# - Encrypts entire startup disk
# - Personal recovery key for admin access
# - Institutional recovery key option available
# - Keys escrowed during initial encryption
# - FIPS 140-2 validated encryption
# - XTS-AES-128 encryption with 256-bit key
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices fully managed by Intune only. Each device must have FileVault encryption enabled with key escrowed to Intune.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to retrieve keys from different types of devices in one action.

**Security Warning:** Retrieved keys will be displayed in action output and may be stored in Terraform state.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

