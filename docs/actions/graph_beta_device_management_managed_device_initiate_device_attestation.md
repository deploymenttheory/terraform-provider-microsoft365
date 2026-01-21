---
page_title: "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates device attestation on managed Windows devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/initiateDeviceAttestation and /deviceManagement/comanagedDevices/{managedDeviceId}/initiateDeviceAttestation endpoints. This action is used to initiate device attestation on managed Windows devices. Device attestation is a security feature that uses the Trusted Platform Module (TPM) to cryptographically verify the device's boot integrity, security configuration, and overall health status. This attestation process creates a trusted baseline that can be used for conditional access, compliance policies, and zero-trust security models. The TPM provides hardware-rooted proof that the device has not been tampered with and is in a known good state.
  Important Notes:
  Only works on Windows devices with TPM 1.2 or 2.0Performs cryptographic verification of device healthCreates attestation report for compliance validationDoes not affect device operation or user accessResults stored in Intune for policy enforcementCritical for Zero Trust security architectureShould be performed periodically for compliance
  Use Cases:
  Conditional access policy enforcementCompliance validation for security standardsZero Trust security model implementationPeriodic device health verificationPre-deployment security validationPost-incident device integrity checks
  Platform Support:
  Windows: Devices with TPM 1.2/2.0 and secure bootOther Platforms: Not supported (Windows-specific feature)
  Reference: Microsoft Graph API - Initiate Device Attestation https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatedeviceattestation?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation (Action)

Initiates device attestation on managed Windows devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/initiateDeviceAttestation` and `/deviceManagement/comanagedDevices/{managedDeviceId}/initiateDeviceAttestation` endpoints. This action is used to initiate device attestation on managed Windows devices. Device attestation is a security feature that uses the Trusted Platform Module (TPM) to cryptographically verify the device's boot integrity, security configuration, and overall health status. This attestation process creates a trusted baseline that can be used for conditional access, compliance policies, and zero-trust security models. The TPM provides hardware-rooted proof that the device has not been tampered with and is in a known good state.

**Important Notes:**
- Only works on Windows devices with TPM 1.2 or 2.0
- Performs cryptographic verification of device health
- Creates attestation report for compliance validation
- Does not affect device operation or user access
- Results stored in Intune for policy enforcement
- Critical for Zero Trust security architecture
- Should be performed periodically for compliance

**Use Cases:**
- Conditional access policy enforcement
- Compliance validation for security standards
- Zero Trust security model implementation
- Periodic device health verification
- Pre-deployment security validation
- Post-incident device integrity checks

**Platform Support:**
- **Windows**: Devices with TPM 1.2/2.0 and secure boot
- **Other Platforms**: Not supported (Windows-specific feature)

**Reference:** [Microsoft Graph API - Initiate Device Attestation](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatedeviceattestation?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [initiateDeviceAttestation action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatedeviceattestation?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Device Attestation and Security Guides
- [Windows device health attestation](https://learn.microsoft.com/en-us/windows/security/operating-system-security/system-security/protect-high-value-assets-by-controlling-the-health-of-windows-10-based-devices)
- [Trusted Platform Module (TPM) overview](https://learn.microsoft.com/en-us/windows/security/hardware-security/tpm/trusted-platform-module-overview)
- [Conditional access device platform condition](https://learn.microsoft.com/en-us/azure/active-directory/conditional-access/concept-conditional-access-conditions#device-platforms)

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
# Example 1: Initiate device attestation on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "initiate_single" {
  config {
    managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Initiate device attestation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "initiate_multiple" {
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
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "initiate_maximal" {
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

# Example 4: Initiate attestation on all Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "initiate_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to initiate device attestation for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to initiate device attestation for. These are devices fully managed by Intune.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to initiate attestation on different types of devices in one action.

**Important:** This action uses the TPM to cryptographically verify device health and security state. It does not affect device operation or user access.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices with TPM before attempting attestation. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

