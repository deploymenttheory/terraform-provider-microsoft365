---
page_title: "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates Mobile Device Management (MDM) key recovery and TPM attestation on managed Windows devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery and /deviceManagement/comanagedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery endpoints. This action is used to perform BitLocker recovery key escrow and Trusted Platform Module (TPM) attestation to ensure recovery keys are properly stored in Azure AD and the device's TPM is healthy. This is critical for security compliance, data recovery scenarios, and ensuring encrypted devices can be recovered if users forget passwords or encounter hardware issues.
  Important Notes:
  Only works on Windows devices with BitLocker and TPM enabledEscrows BitLocker recovery keys to Azure ADPerforms TPM health attestationDoes not encrypt/decrypt the deviceDoes not affect device operation or user accessEssential for compliance and disaster recoveryShould be run periodically for key rotation
  Use Cases:
  Ensuring BitLocker recovery keys are escrowedCompliance auditing for encryption key managementVerifying TPM attestation and healthPeriodic key rotation and refreshPre-deployment validation for new devicesRecovery preparation for critical devices
  Platform Support:
  Windows: Devices with BitLocker and TPM 1.2/2.0Other Platforms: Not supported (Windows-specific feature)
  Reference: Microsoft Graph API - Initiate Mobile Device Management Key Recovery https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatemobiledevicemanagementkeyrecovery?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery (Action)

Initiates Mobile Device Management (MDM) key recovery and TPM attestation on managed Windows devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` and `/deviceManagement/comanagedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` endpoints. This action is used to perform BitLocker recovery key escrow and Trusted Platform Module (TPM) attestation to ensure recovery keys are properly stored in Azure AD and the device's TPM is healthy. This is critical for security compliance, data recovery scenarios, and ensuring encrypted devices can be recovered if users forget passwords or encounter hardware issues.

**Important Notes:**
- Only works on Windows devices with BitLocker and TPM enabled
- Escrows BitLocker recovery keys to Azure AD
- Performs TPM health attestation
- Does not encrypt/decrypt the device
- Does not affect device operation or user access
- Essential for compliance and disaster recovery
- Should be run periodically for key rotation

**Use Cases:**
- Ensuring BitLocker recovery keys are escrowed
- Compliance auditing for encryption key management
- Verifying TPM attestation and health
- Periodic key rotation and refresh
- Pre-deployment validation for new devices
- Recovery preparation for critical devices

**Platform Support:**
- **Windows**: Devices with BitLocker and TPM 1.2/2.0
- **Other Platforms**: Not supported (Windows-specific feature)

**Reference:** [Microsoft Graph API - Initiate Mobile Device Management Key Recovery](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatemobiledevicemanagementkeyrecovery?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [initiateMobileDeviceManagementKeyRecovery action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatemobiledevicemanagementkeyrecovery?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### BitLocker and TPM Management Guides
- [BitLocker recovery keys in Azure AD](https://learn.microsoft.com/en-us/azure/active-directory/devices/device-management-azure-portal#bitlocker-recovery-keys)
- [Manage BitLocker policy for Windows devices with Intune](https://learn.microsoft.com/en-us/mem/intune/protect/encrypt-devices)
- [Trusted Platform Module (TPM) technology overview](https://learn.microsoft.com/en-us/windows/security/hardware-security/tpm/trusted-platform-module-overview)

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
# Example 1: Initiate MDM key recovery on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_single" {
  config {
    managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Initiate MDM key recovery on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_multiple" {
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

# Example 3: Initiate with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_maximal" {
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

# Example 4: Initiate key recovery on all iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "initiate_all_ios" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. These are devices fully managed by Intune.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to initiate key recovery on different types of devices in one action.

**Important:** This action escrows BitLocker recovery keys to Azure AD and performs TPM attestation. It does not affect device operation or user access.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting key recovery. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

