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
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


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
# Example 1: Rotate FileVault key on a single macOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Rotate FileVault keys on multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_multiple" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Rotate FileVault keys with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Rotate FileVault keys on all macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_all_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 5: Rotate FileVault keys for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_comanaged" {
  config {
    comanaged_device_ids = [
      "11111111-1111-1111-1111-111111111111",
      "22222222-2222-2222-2222-222222222222"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Rotate FileVault keys for company-owned macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "company_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_company_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.company_macos.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "rotated_filevault_keys_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key.rotate_multiple.config.managed_device_ids)
  description = "Number of devices that had FileVault keys rotated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to rotate FileVault keys for. These are macOS devices fully managed by Intune only. Each device must have FileVault encryption enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to rotate keys on different types of devices in one action.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are macOS devices before attempting to rotate FileVault keys. Disabling this can speed up planning but may result in runtime errors for non-existent or non-macOS devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

