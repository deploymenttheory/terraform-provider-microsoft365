---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password Action - terraform-provider-microsoft365"
subcategory: "Actions - Device Management"
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

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_rotate_local_admin_password (Action)

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

## Example Usage

### Basic - Managed Devices

```terraform
# Example 1: Basic - Rotate local admin password for managed devices
# Use case: Regular security password rotation (quarterly security maintenance)
action "rotate_password_basic" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc", # Windows 10 workstation
    "87654321-4321-4321-4321-ba9876543210", # Windows 11 laptop
  ]

  timeouts {
    invoke = "10m"
  }
}

# Example 2: Co-managed devices - Hybrid SCCM and Intune environment
# Use case: Password rotation in hybrid cloud/on-prem environment
action "rotate_password_comanaged" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  comanaged_device_ids = [
    "abcdef12-3456-7890-abcd-ef1234567890", # SCCM + Intune device
  ]

  timeouts {
    invoke = "5m"
  }
}

# Example 3: Mixed environment - Both managed and co-managed devices
# Use case: Large-scale password rotation across entire Windows fleet
action "rotate_password_mixed" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc", # Pure Intune devices
    "22334455-6677-8899-aabb-ccddeefffabc",
  ]

  comanaged_device_ids = [
    "abcdef12-3456-7890-abcd-ef1234567890", # Hybrid SCCM devices
    "fedcba09-8765-4321-fedc-ba0987654321",
  ]

  timeouts {
    invoke = "15m"
  }
}

# Example 4: Post-incident response - Emergency password rotation
# Use case: Security incident requiring immediate password reset
action "rotate_password_emergency" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    "compromised-device-guid-1", # Device with potential compromise
    "compromised-device-guid-2", # Related device in same network
    "compromised-device-guid-3", # Device accessed by same admin
  ]

  timeouts {
    invoke = "10m"
  }
}

# Example 5: Compliance-driven rotation - Quarterly security compliance
# Use case: Scheduled password rotation for PCI-DSS, HIPAA, or other compliance
action "rotate_password_compliance_q1" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    # Finance department devices (PCI-DSS compliance)
    "finance-device-1",
    "finance-device-2",
    "finance-device-3",
    # Healthcare devices (HIPAA compliance)
    "medical-workstation-1",
    "medical-workstation-2",
  ]

  timeouts {
    invoke = "15m"
  }
}

# Example 6: Admin lifecycle - Offboarding administrator
# Use case: Administrator leaving organization, rotate all devices they had access to
action "rotate_password_admin_offboarding" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    "it-admin-workstation-1",  # Admin's primary workstation
    "server-mgmt-device-1",    # Server management device
    "backup-admin-device-1",   # Backup admin access device
  ]

  timeouts {
    invoke = "10m"
  }
}

# Example 7: Zero Trust implementation - Rotate passwords for privileged access workstations (PAWs)
# Use case: Zero Trust security model requiring frequent password rotation
action "rotate_password_zero_trust" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    "paw-tier0-1", # Domain Admin PAW
    "paw-tier0-2", # Enterprise Admin PAW
    "paw-tier1-1", # Server Admin PAW
    "paw-tier1-2", # Application Admin PAW
    "paw-tier2-1", # Workstation Admin PAW
  ]

  timeouts {
    invoke = "15m"
  }
}

# Example 8: Audit preparation - Pre-audit password validation
# Use case: Preparing for security audit, ensure all passwords are recently rotated
action "rotate_password_pre_audit" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    # Critical infrastructure devices
    "domain-controller-mgmt",
    "exchange-admin-device",
    "sql-admin-device",
    "backup-mgmt-device",
  ]

  comanaged_device_ids = [
    # Hybrid devices requiring audit compliance
    "sccm-admin-device",
    "hybrid-server-mgmt",
  ]

  timeouts {
    invoke = "10m"
  }
}

# Example 9: Using data sources to dynamically select devices
# Use case: Rotate passwords for all Windows devices in specific group
data "microsoft365_graph_beta_device_management_managed_device_list" "windows_devices" {
  filter = "operatingSystem eq 'Windows'"
}

action "rotate_password_dynamic" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device_list.windows_devices.managed_devices : device.id]

  timeouts {
    invoke = "30m"
  }
}

# Example 10: Departmental rotation - IT department quarterly maintenance
# Use case: Regular password rotation for IT department devices
action "rotate_password_it_dept" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_rotate_local_admin_password

  managed_device_ids = [
    # Helpdesk devices
    "helpdesk-1",
    "helpdesk-2",
    # IT admin devices
    "it-admin-1",
    "it-admin-2",
    # Network admin devices
    "network-admin-1",
    # Security admin devices
    "security-admin-1",
  ]

  timeouts {
    invoke = "15m"
  }
}

# Important Notes:
#
# 1. Windows LAPS Requirements:
#    - Devices must be Windows 10/11
#    - Windows LAPS policy must be configured in Intune
#    - Policy must be assigned and active on target devices
#    - Devices must have checked in and applied the policy
#
# 2. Password Generation:
#    - Passwords automatically generated (complex, random)
#    - Length and complexity based on LAPS policy settings
#    - Stored securely in Azure AD or Intune
#    - Previous password immediately invalidated
#
# 3. Password Retrieval:
#    - Authorized administrators can retrieve passwords
#    - Access requires appropriate Azure AD permissions
#    - Retrieval actions are audited
#    - Passwords can be retrieved via Azure Portal or Graph API
#
# 4. Security Best Practices:
#    - Rotate passwords regularly (quarterly or semi-annually)
#    - Rotate immediately after admin offboarding
#    - Rotate after suspected compromise
#    - Rotate as part of incident response
#    - Document rotation schedule for compliance
#
# 5. Zero Trust Considerations:
#    - Align rotation with Zero Trust principles
#    - Use just-in-time (JIT) password retrieval
#    - Enable MFA for password retrieval
#    - Monitor password retrieval events
#    - Regular rotation reduces exposure window
#
# 6. Compliance Requirements:
#    - Many frameworks require regular password rotation
#    - NIST recommends rotation after compromise
#    - PCI-DSS requires periodic password changes
#    - HIPAA security rule requires password management
#    - SOC 2 Type II requires privileged access controls
```

## API Permissions

The following [Microsoft Graph API permissions](https://learn.microsoft.com/en-us/graph/permissions-reference) are required to use this action:

### Read Permissions
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

### Write Permissions
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

## Authentication Mechanism

This action uses the configured Microsoft 365 provider authentication. Ensure you have properly configured the provider with appropriate credentials and permissions.

## Important Notes

### Windows LAPS Requirements
- Devices must be Windows 10 or Windows 11
- Windows LAPS policy must be configured and assigned in Intune
- Policy must be active and applied on target devices
- Devices must support Windows LAPS (modern Windows versions)

### Password Generation and Storage
- Passwords are automatically generated as complex, random strings
- Length and complexity based on LAPS policy settings
- Stored securely in Azure AD or Intune
- Previous password immediately invalidated upon rotation
- Password history maintained for audit purposes

### Password Retrieval
- Only authorized administrators can retrieve passwords
- Retrieval requires specific Azure AD role permissions:
  - Cloud Device Administrator
  - Intune Administrator
  - Global Administrator
- Password retrieval actions are audited in Azure AD logs
- Passwords can be retrieved via:
  - Azure Portal (Azure AD > Devices)
  - Microsoft Graph API
  - Intune Admin Center

### Security Best Practices
1. **Regular Rotation Schedule**: Implement quarterly or semi-annual rotation
2. **Incident Response**: Rotate immediately after:
   - Administrator offboarding
   - Suspected credential compromise
   - Security incidents
   - Failed audit findings
3. **Access Control**: 
   - Limit password retrieval permissions
   - Enable MFA for password retrieval
   - Monitor retrieval events
4. **Zero Trust Alignment**:
   - Use just-in-time (JIT) password retrieval
   - Implement time-limited access
   - Regular rotation reduces exposure window

### Compliance Considerations
This action helps meet compliance requirements for:
- **PCI-DSS**: Privileged account password management
- **HIPAA**: Security rule for password controls
- **SOC 2 Type II**: Privileged access controls
- **NIST**: Password management and rotation after compromise
- **ISO 27001**: Information security controls for privileged access

### Platform Support
- **Windows 10/11**: Fully supported with Windows LAPS
- **Earlier Windows versions**: Not supported
- **Non-Windows platforms**: Not applicable (LAPS is Windows-specific)

### Operation Impact
- No device restart required
- No user interruption
- No service disruption
- Immediate effect on local administrator account
- Previous password becomes invalid immediately

## Microsoft Graph API References

- [rotateLocalAdminPassword action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatelocaladminpassword?view=graph-rest-beta)
- [Windows LAPS in Intune](https://learn.microsoft.com/en-us/mem/intune/protect/windows-laps-overview)
- [Local Administrator Password Solution (LAPS)](https://learn.microsoft.com/en-us/windows-server/identity/laps/laps-overview)

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to rotate local administrator passwords for. These are devices managed by both Intune and Configuration Manager (SCCM) with Windows LAPS enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to rotate local administrator passwords for. These are devices fully managed by Intune with Windows LAPS enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to rotate passwords on different types of devices in one action.

**Important:** Devices must have Windows LAPS policy configured and enabled. The new password will be automatically generated and stored securely in Azure AD or Intune.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

