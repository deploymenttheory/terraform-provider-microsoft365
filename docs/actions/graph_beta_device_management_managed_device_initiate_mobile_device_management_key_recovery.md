---
page_title: "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates Mobile Device Management (MDM) key recovery and TPM attestation on managed Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery and /deviceManagement/comanagedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery endpoints. This action performs BitLocker recovery key escrow and Trusted Platform Module (TPM) attestation to ensure recovery keys are properly stored in Azure AD and the device's TPM is healthy. This is critical for security compliance, data recovery scenarios, and ensuring encrypted devices can be recovered if users forget passwords or encounter hardware issues.
  Important Notes:
  Only works on Windows devices with BitLocker and TPM enabledEscrows BitLocker recovery keys to Azure ADPerforms TPM health attestationDoes not encrypt/decrypt the deviceDoes not affect device operation or user accessEssential for compliance and disaster recoveryShould be run periodically for key rotation
  Use Cases:
  Ensuring BitLocker recovery keys are escrowedCompliance auditing for encryption key managementVerifying TPM attestation and healthPeriodic key rotation and refreshPre-deployment validation for new devicesRecovery preparation for critical devices
  Platform Support:
  Windows: Devices with BitLocker and TPM 1.2/2.0Other Platforms: Not supported (Windows-specific feature)
  Reference: Microsoft Graph API - Initiate Mobile Device Management Key Recovery https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiatemobiledevicemanagementkeyrecovery?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery (Action)

Initiates Mobile Device Management (MDM) key recovery and TPM attestation on managed Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` and `/deviceManagement/comanagedDevices/{managedDeviceId}/initiateMobileDeviceManagementKeyRecovery` endpoints. This action performs BitLocker recovery key escrow and Trusted Platform Module (TPM) attestation to ensure recovery keys are properly stored in Azure AD and the device's TPM is healthy. This is critical for security compliance, data recovery scenarios, and ensuring encrypted devices can be recovered if users forget passwords or encounter hardware issues.

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
| **Windows** | ✅ Full Support | BitLocker enabled, TPM 1.2 or 2.0 present |
| **macOS** | ❌ Not Supported | FileVault keys managed separately |
| **iOS/iPadOS** | ❌ Not Supported | No BitLocker equivalent |
| **Android** | ❌ Not Supported | Encryption managed differently |

### What is MDM Key Recovery?

MDM (Mobile Device Management) key recovery is a process that:
- Escrows BitLocker recovery keys to Azure Active Directory
- Performs Trusted Platform Module (TPM) attestation
- Verifies device encryption and security posture
- Ensures recovery keys are accessible for disaster recovery
- Does not encrypt or decrypt the device
- Does not affect device operation or user access
- Is critical for compliance and data protection

### What Happens During Key Recovery

1. **Key Escrow**: Device uploads BitLocker recovery keys to Azure AD
2. **TPM Attestation**: Device firmware (TPM) validates its security state
3. **Azure AD Storage**: Recovery keys are securely stored in Azure AD
4. **Audit Logging**: Action is logged for compliance tracking
5. **Completion**: Device reports success without user impact

### BitLocker Recovery Keys

**What are BitLocker Recovery Keys?**
- 48-digit numerical passwords used to unlock encrypted drives
- Generated when BitLocker is enabled on a Windows device
- Required if user forgets password or BitLocker suspects tampering
- Critical for data recovery in disaster scenarios

**Where are Keys Stored?**
- Azure Active Directory (via this action)
- Local Active Directory (for domain-joined devices)
- User's Microsoft account (for personal devices)
- Printed or saved by user during setup

**Accessing Recovery Keys:**
- Azure AD Portal → Devices → Device Details → BitLocker Keys
- Microsoft Graph API queries
- PowerShell with Azure AD module
- End-user self-service portal (if enabled)

### TPM Attestation

**What is TPM?**
- **Trusted Platform Module**: Hardware security chip on motherboard
- Stores cryptographic keys, passwords, and certificates
- Provides hardware-based security functions
- Required for modern Windows security features

**TPM Versions:**
| Version | Year | Features |
|---------|------|----------|
| **TPM 1.2** | 2005 | Basic cryptographic operations, BitLocker support |
| **TPM 2.0** | 2014 | Enhanced algorithms, better performance, Windows 11 required |

**Attestation Process:**
- TPM provides cryptographic proof of device health
- Validates boot integrity and security configuration
- Confirms no tampering with device firmware
- Reports measured boot process to cloud service

### When to Initiate Key Recovery

- **New Device Enrollment**: Ensure keys are escrowed immediately
- **Compliance Audits**: Verify all encrypted devices have keys backed up
- **Periodic Rotation**: Refresh keys quarterly or semi-annually
- **Pre-Deployment**: Validate encryption before user assignment
- **After BitLocker Changes**: Re-escrow keys after policy updates
- **Troubleshooting**: Verify key escrow if recovery issues occur
- **Security Incidents**: Confirm TPM attestation after suspected compromise

### Security and Compliance

**Compliance Requirements:**
- Many regulations require recovery key escrow (GDPR, HIPAA, PCI-DSS)
- Keys must be stored securely and access-controlled
- Audit logs must track key access and recovery operations
- Regular key rotation recommended for security

**Security Benefits:**
- **Data Recovery**: Access encrypted data if user locked out
- **Device Wipe Assurance**: Verify encryption before device disposal
- **Compliance Proof**: Demonstrate encryption key management
- **Reduced Downtime**: Quick recovery from BitLocker lockout
- **TPM Validation**: Confirm device hasn't been tampered with

### Key Recovery vs Other Actions

| Action | Purpose | Changes Device |
|--------|---------|----------------|
| **MDM Key Recovery** | Escrow keys + TPM check | ❌ No |
| **Rotate BitLocker Keys** | Generate new encryption keys | ✅ Yes |
| **Get FileVault Key** | Retrieve macOS recovery key | ❌ No |
| **Wipe Device** | Factory reset with encryption | ✅ Yes |

### Important Considerations

✅ **Safe Operations:**
- Does not affect device performance
- No user downtime or interruption
- No data encryption/decryption
- Can be run anytime, on any device
- Idempotent (safe to run multiple times)

⚠️ **Requirements:**
- Device must have BitLocker enabled
- TPM must be present and functional
- Device must be online and checking in
- Adequate Azure AD storage for keys
- Proper permissions configured

### Troubleshooting

**Common Issues:**

1. **Device Not Encrypted**
   - Solution: Enable BitLocker via Intune policy first
   - Check: Device must show as encrypted in Intune

2. **No TPM Present**
   - Solution: Cannot proceed, device lacks hardware
   - Alternative: Use software-only encryption (less secure)

3. **Key Not Appearing in Azure AD**
   - Solution: Wait 15-30 minutes for sync
   - Check: Device last check-in time and connectivity

4. **Action Fails Silently**
   - Solution: Check Intune logs and device compliance
   - Verify: Device has proper encryption policy applied

### Best Practices

**Operational:**
- ✅ Escrow keys during device enrollment
- ✅ Schedule periodic key refresh (quarterly)
- ✅ Validate key escrow after BitLocker policy changes
- ✅ Monitor key escrow success rates
- ✅ Document key access procedures
- ✅ Test key retrieval process regularly

**Security:**
- ✅ Restrict access to recovery keys (RBAC)
- ✅ Audit all key access and retrieval
- ✅ Rotate keys periodically for sensitive devices
- ✅ Use MFA for key access
- ✅ Monitor for TPM attestation failures
- ✅ Investigate failed attestations promptly

**Compliance:**
- ✅ Document key management procedures
- ✅ Maintain audit trail of all key operations
- ✅ Regular compliance validation
- ✅ Include in disaster recovery plans
- ✅ Train helpdesk on key recovery procedures

## Example Usage

```terraform
# Example 1: Initiate MDM key recovery for single device
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "single_device" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Initiate key recovery for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "multiple_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Compliance-driven key recovery
variable "compliance_devices" {
  description = "Device IDs requiring BitLocker key escrow for compliance"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "compliance_escrow" {
  managed_device_ids = var.compliance_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Periodic key rotation based on data source
data "microsoft365_graph_beta_device_management_managed_device" "encrypted_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceCategoryDisplayName eq 'Encrypted'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "periodic_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.encrypted_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: New device enrollment key escrow
locals {
  new_enrollment_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "new_device_escrow" {
  managed_device_ids = local.new_enrollment_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Co-managed device key recovery
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "comanaged_escrow" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Critical devices pre-deployment validation
data "microsoft365_graph_beta_device_management_managed_device" "critical_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Critical Infrastructure'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "critical_validation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.critical_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Department-specific key recovery
locals {
  finance_department_devices = {
    "finance_laptop_1" = "11111111-1111-1111-1111-111111111111"
    "finance_laptop_2" = "22222222-2222-2222-2222-222222222222"
    "finance_desktop_1" = "33333333-3333-3333-3333-333333333333"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "finance_department" {
  managed_device_ids = values(local.finance_department_devices)

  timeouts = {
    invoke = "15m"
  }
}

# Example 9: Scheduled quarterly key rotation
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "quarterly_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : device.id]

  timeouts = {
    invoke = "60m"
  }
}

# Example 10: Audit preparation key escrow
locals {
  audit_scope_devices = [
    "audit01-1111-1111-1111-111111111111",
    "audit02-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "audit_preparation" {
  managed_device_ids = local.audit_scope_devices

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "key_recovery_summary" {
  value = {
    managed   = length(action.multiple_devices.managed_device_ids)
    comanaged = length(action.comanaged_escrow.comanaged_device_ids)
  }
  description = "Count of devices with MDM key recovery initiated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to initiate MDM key recovery and TPM attestation for. These are devices fully managed by Intune.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to initiate key recovery on different types of devices in one action.

**Important:** This action escrows BitLocker recovery keys to Azure AD and performs TPM attestation. It does not affect device operation or user access.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

