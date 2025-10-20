# Example 1: Rotate FileVault key for a single macOS device
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Rotate FileVault keys for multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_multiple" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Rotate keys for all managed macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_all_macos" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 4: Rotate keys for macOS devices in specific department
data "microsoft365_graph_beta_device_management_managed_device" "finance_macos" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS' and deviceCategoryDisplayName eq 'Finance'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_finance_macos" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_macos.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Scheduled quarterly key rotation
locals {
  # Devices due for quarterly rotation
  macos_devices_for_rotation = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "quarterly_rotation" {
  managed_device_ids = local.macos_devices_for_rotation

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Rotate key for co-managed macOS device
action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_comanaged" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Rotate keys after security incident
variable "compromised_device_ids" {
  description = "List of device IDs potentially affected by security incident"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "incident_response" {
  managed_device_ids = var.compromised_device_ids

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Rotate keys for devices being reassigned
data "microsoft365_graph_beta_device_management_managed_device" "reassignment_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'previous.user@example.com' and operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "reassignment_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.reassignment_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 9: Rotate keys for devices with old enrollment dates
data "microsoft365_graph_beta_device_management_managed_device" "old_macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS' and enrolledDateTime lt 2024-01-01T00:00:00Z"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_file_vault_key" "rotate_old_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_macos_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Output examples
output "devices_rotated_count" {
  value       = length(action.rotate_multiple.managed_device_ids)
  description = "Number of devices that had FileVault keys rotated"
}

output "rotation_summary" {
  value = {
    managed   = length(action.rotate_all_macos.managed_device_ids)
    comanaged = length(action.rotate_comanaged.comanaged_device_ids)
  }
  description = "Count of FileVault key rotations by device type"
}

# Important Notes:
# FileVault Key Rotation Features:
# - Generates new FileVault recovery key
# - New key is automatically escrowed with Intune
# - Previous recovery key becomes invalid
# - No user interaction required
# - Device does not need to be logged in
# - Works silently in the background
# - Requires macOS devices with FileVault enabled
#
# When to Rotate FileVault Keys:
# - Regular compliance-driven key rotation (quarterly/annually)
# - After suspected key compromise or exposure
# - When changing device ownership or assignment
# - As part of security incident response
# - Before device reassignment to new users
# - After employee termination or transfer
# - Compliance requirements mandate key rotation
#
# What Happens When Key is Rotated:
# - Device receives rotation command from Intune
# - macOS generates new FileVault recovery key
# - New key is escrowed with Intune automatically
# - Previous recovery key is invalidated
# - Change occurs without user interaction
# - No device restart required
# - User passwords remain unchanged
# - Disk encryption continues uninterrupted
#
# Platform Requirements:
# - macOS: Fully supported (FileVault-enabled devices)
# - Device must have FileVault encryption enabled
# - Device must be enrolled in Intune
# - Device must be online to receive command
# - Other platforms: Not applicable (FileVault is macOS-only)
#
# Best Practices:
# - Implement regular rotation schedule (quarterly)
# - Rotate keys after employee changes
# - Document rotation policy and schedule
# - Test rotation on pilot devices first
# - Monitor rotation success/failure rates
# - Keep audit logs of all rotations
# - Rotate immediately after suspected compromise
# - Combine with other security measures
#
# Security Benefits:
# - Limits exposure window if key compromised
# - Compliance with security policies
# - Part of defense-in-depth strategy
# - Reduces risk of unauthorized access
# - Maintains control over device recovery
# - Supports zero-trust security model
# - Helps meet regulatory requirements
#
# FileVault Overview:
# - Full disk encryption for macOS
# - Encrypts entire startup disk
# - Recovery key allows admin unlock
# - Required for many compliance standards
# - Protects data at rest
# - Transparent to users
# - Minimal performance impact
#
# Recovery Key Management:
# - Keys escrowed with Intune
# - Accessible to admins via Intune portal
# - Used when user forgets password
# - Required for device recovery
# - Rotation creates new unique key
# - Old key no longer works after rotation
# - Keys stored securely in Intune
#
# Compliance Considerations:
# - Many standards require key rotation
# - NIST recommends periodic rotation
# - Document rotation frequency
# - Maintain rotation audit trail
# - Verify rotation completion
# - Report on rotation compliance
# - Integrate with compliance dashboards
#
# Troubleshooting:
# - Verify device is macOS
# - Check FileVault is enabled
# - Ensure device is online
# - Verify Intune connectivity
# - Check device compliance status
# - Review device logs if rotation fails
# - Contact Microsoft support if needed
#
# Related Actions:
# - Device compliance: Enforce FileVault
# - Encryption policies: Configure FileVault
# - Retire device: Remove escrowed keys
# - Remote lock: Use with rotated keys
# - Device wipe: Secure data removal
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rotatefilevaultkey?view=graph-rest-beta

