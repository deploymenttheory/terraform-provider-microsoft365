---
page_title: "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves the FileVault recovery key for macOS managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/getFileVaultKey and /deviceManagement/comanagedDevices/{managedDeviceId}/getFileVaultKey endpoints. This action is used to retrieve escrowed FileVault recovery keys for device recovery purposes. The recovery key is displayed in the action output and can be used to unlock an encrypted macOS device when a user has forgotten their password or is otherwise unable to access the device. This is a critical capability for IT support and device recovery scenarios.
  Important Security Notes:
  Recovery keys are highly sensitive credentialsKeys grant full access to encrypted device dataAccess to keys should be audited and restrictedOnly retrieve keys when necessary for device recoveryKeys are displayed in plain text in action outputEnsure proper security controls on Terraform stateConsider security implications before using in automation
  Use Cases:
  Emergency device recovery when user cannot log inUnlocking devices for departing employeesTechnical support scenarios requiring device accessDisaster recovery and business continuityDevice repurposing or reassignment preparation
  Platform Support:
  macOS: Fully supported on devices with FileVault enabled and keys escrowedOther Platforms: Not applicable (FileVault is macOS-only)
  Reference: Microsoft Graph API - Get FileVault Key https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_get_file_vault_key (Action)

Retrieves the FileVault recovery key for macOS managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/getFileVaultKey` and `/deviceManagement/comanagedDevices/{managedDeviceId}/getFileVaultKey` endpoints. This action is used to retrieve escrowed FileVault recovery keys for device recovery purposes. The recovery key is displayed in the action output and can be used to unlock an encrypted macOS device when a user has forgotten their password or is otherwise unable to access the device. This is a critical capability for IT support and device recovery scenarios.

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

## Microsoft Documentation

### Graph API References
- [getFileVaultKey function](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-getfilevaultkey?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune FileVault Guides
- [Use disk encryption for macOS with Intune](https://learn.microsoft.com/en-us/mem/intune/protect/encrypt-devices-filevault)
- [FileVault recovery key rotation](https://learn.microsoft.com/en-us/mem/intune/protect/encrypt-devices-filevault#rotate-recovery-keys)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Example Usage

```terraform
# Example 1: Retrieve FileVault key for a single macOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Retrieve FileVault keys for multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_multiple" {
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

# Example 3: Retrieve keys with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_with_validation" {
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

# Example 4: Emergency recovery for locked macOS devices
variable "locked_device_ids" {
  description = "Device IDs for locked macOS devices"
  type        = list(string)
  default = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "emergency_recovery" {
  config {
    managed_device_ids = var.locked_device_ids

    validate_device_exists = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Retrieve keys for departing employee's macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_employee" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'departing.employee@example.com') and (operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "departing_employee_recovery" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_employee.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Retrieve keys for co-managed macOS devices
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_comanaged" {
  config {
    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Retrieve keys for all company macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "all_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "all_company_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_macos.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "retrieved_keys_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_get_file_vault_key.retrieve_multiple.config.managed_device_ids)
  description = "Number of devices for which FileVault keys were retrieved"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to retrieve FileVault keys for. These are macOS devices fully managed by Intune only. Each device must have FileVault encryption enabled with key escrowed to Intune.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to retrieve keys from different types of devices in one action.

**Security Warning:** Retrieved keys will be displayed in action output and may be stored in Terraform state.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are macOS devices with FileVault enabled before attempting to retrieve keys. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

