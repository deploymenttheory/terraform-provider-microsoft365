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
    "it-admin-workstation-1", # Admin's primary workstation
    "server-mgmt-device-1",   # Server management device
    "backup-admin-device-1",  # Backup admin access device
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

