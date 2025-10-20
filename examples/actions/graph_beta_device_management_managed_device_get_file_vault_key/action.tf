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

