---
page_title: "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates manual rotation of the local administrator password on managed Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/rotateLocalAdminPassword and /deviceManagement/comanagedDevices/{managedDeviceId}/rotateLocalAdminPassword endpoints. This action works with Windows Local Administrator Password Solution (LAPS) to generate and rotate local admin passwords on Windows devices. The new password is automatically generated, stored securely in Azure AD or Intune, and can be retrieved by authorized administrators. This enhances security by ensuring regular password rotation and centralized password management for local administrator accounts.
  Important Notes:
  Only works on Windows 10/11 devices with Windows LAPS enabledRequires Windows LAPS policy configured in IntuneNew password automatically generated and stored in Azure AD/IntunePassword retrievable by authorized administratorsDoes not affect device operation or require restartCritical for security compliance and privileged access managementShould be performed regularly or after admin account compromise
  Use Cases:
  Regular security password rotation (quarterly/semi-annually)Post-incident password reset after compromiseCompliance requirements for privileged account managementOnboarding/offboarding administrator accessAudit preparation and security validationZero Trust privileged access implementation
  Platform Support:
  Windows: Windows 10/11 with Windows LAPS enabledOther Platforms: Not supported (Windows LAPS-specific)
  Reference: Microsoft Graph API - Rotate Local Admin Password https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatelocaladminpassword?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password (Action)

Initiates manual rotation of the local administrator password on managed Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/rotateLocalAdminPassword` and `/deviceManagement/comanagedDevices/{managedDeviceId}/rotateLocalAdminPassword` endpoints. This action works with Windows Local Administrator Password Solution (LAPS) to generate and rotate local admin passwords on Windows devices. The new password is automatically generated, stored securely in Azure AD or Intune, and can be retrieved by authorized administrators. This enhances security by ensuring regular password rotation and centralized password management for local administrator accounts.

**Important Notes:**
- Only works on Windows 10/11 devices with Windows LAPS enabled
- Requires Windows LAPS policy configured in Intune
- New password automatically generated and stored in Azure AD/Intune
- Password retrievable by authorized administrators
- Does not affect device operation or require restart
- Critical for security compliance and privileged access management
- Should be performed regularly or after admin account compromise

**Use Cases:**
- Regular security password rotation (quarterly/semi-annually)
- Post-incident password reset after compromise
- Compliance requirements for privileged account management
- Onboarding/offboarding administrator access
- Audit preparation and security validation
- Zero Trust privileged access implementation

**Platform Support:**
- **Windows**: Windows 10/11 with Windows LAPS enabled
- **Other Platforms**: Not supported (Windows LAPS-specific)

**Reference:** [Microsoft Graph API - Rotate Local Admin Password](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatelocaladminpassword?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [rotateLocalAdminPassword action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatelocaladminpassword?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune LAPS Guides
- [Windows LAPS in Intune](https://learn.microsoft.com/en-us/mem/intune/protect/windows-laps-overview)
- [Local Administrator Password Solution (LAPS)](https://learn.microsoft.com/en-us/windows-server/identity/laps/laps-overview)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`
- **Delegated**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows 10** | ✅ Full Support | Windows LAPS policy configured and assigned |
| **Windows 11** | ✅ Full Support | Windows LAPS policy configured and assigned |
| **Windows 8.1 and earlier** | ❌ Not Supported | Windows LAPS not available |
| **macOS** | ❌ Not Supported | LAPS is Windows-specific |
| **iOS/iPadOS** | ❌ Not Supported | LAPS is Windows-specific |
| **Android** | ❌ Not Supported | LAPS is Windows-specific |

### What is Windows LAPS Password Rotation?

Windows LAPS Password Rotation is an action that:
- Generates a new complex, random local administrator password
- Automatically stores the new password securely in Azure AD or Intune
- Invalidates the previous local administrator password
- Operates without user interaction or device restart
- Enhances security through regular credential changes
- Maintains audit trail of password changes

### When to Rotate Local Admin Passwords

- Regular compliance-driven rotation (quarterly/semi-annually per security policy)
- After suspected password compromise or exposure
- When reassigning devices to new users or departments
- As part of security incident response procedures
- Before or after employee termination or transfer (IT admin access)
- To meet regulatory or audit requirements (PCI-DSS, HIPAA, SOC 2)
- After password has been accessed by administrative staff
- Following privilege escalation security incidents
- As required by Zero Trust security model

### What Happens When Password is Rotated

- Intune sends rotation command to the Windows device
- Device generates new complex password per LAPS policy settings
- New password is automatically escrowed with Azure AD/Intune
- Previous local administrator password is invalidated immediately
- Process completes without user interaction or awareness
- No device restart or user interruption required
- Password history is maintained for audit purposes
- New password becomes available in Intune portal for authorized admin retrieval

### Password Retrieval and Access

**Who Can Retrieve Passwords:**
- Only authorized administrators with specific Azure AD role permissions:
  - Cloud Device Administrator
  - Intune Administrator
  - Global Administrator

**How to Retrieve Passwords:**
- Azure Portal (Azure AD > Devices)
- Microsoft Graph API
- Intune Admin Center

**Security Measures:**
- Password retrieval actions are audited in Azure AD logs
- Passwords stored securely encrypted in Azure AD/Intune
- Just-in-time (JIT) retrieval model recommended
- Multi-factor authentication required for retrieval
- Time-limited access windows enforced

### Security Best Practices

- Implement regular rotation schedule (quarterly or semi-annual)
- Rotate immediately after administrator offboarding
- Rotate immediately after suspected credential compromise
- Limit password retrieval permissions to essential staff
- Enable MFA for all password retrieval operations
- Monitor and alert on password retrieval events
- Use just-in-time (JIT) password access patterns
- Implement time-limited access windows
- Regular rotation reduces credential exposure window
- Maintain audit logs for compliance requirements

### Compliance Considerations

This action helps meet compliance requirements for:

| Framework | Requirement |
|-----------|-------------|
| **PCI-DSS** | Privileged account password management and rotation |
| **HIPAA** | Security rule for password controls and privileged access |
| **SOC 2 Type II** | Privileged access controls and credential management |
| **NIST** | Password management and rotation after compromise |
| **ISO 27001** | Information security controls for privileged access |

### Common Use Cases

- Scheduled compliance-driven password rotation
- Emergency password rotation after security incidents
- Device provisioning and deprovisioning workflows
- IT administrator offboarding procedures
- Privilege access management (PAM) integration
- Zero Trust security model implementation
- Audit preparation and remediation
- Security incident response and recovery
- Device reassignment to different users/departments
- Regular security hygiene maintenance

## Example Usage

```terraform
# Example 1: Rotate local admin password on a single Windows device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Rotate local admin passwords on multiple Windows devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_multiple" {
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

# Example 3: Rotate local admin passwords with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_with_validation" {
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

# Example 4: Rotate local admin passwords on all Windows 10/11 devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 5: Rotate local admin passwords for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "rotate_comanaged" {
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

# Example 6: Scheduled rotation for compliance
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password" "scheduled_rotation" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

    ignore_partial_failures = false

    timeouts = {
      invoke = "25m"
    }
  }
}

# Output examples
output "rotated_passwords_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password.rotate_multiple.config.managed_device_ids)
  description = "Number of devices that had local admin passwords rotated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to rotate local administrator passwords for. These are devices managed by both Intune and Configuration Manager (SCCM) with Windows LAPS enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to rotate local administrator passwords for. These are devices fully managed by Intune with Windows LAPS enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to rotate passwords on different types of devices in one action.

**Important:** Devices must have Windows LAPS policy configured and enabled. The new password will be automatically generated and stored securely in Azure AD or Intune.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting to rotate local admin passwords. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

