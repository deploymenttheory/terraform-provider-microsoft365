# ============================================================================
# Example 1: Recover passcode for a single iOS device by ID
# ============================================================================
# Use case: User forgot passcode on supervised iPad
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "single_device" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 2: Recover passcodes for multiple supervised iOS devices
# ============================================================================
# Use case: Help desk batch recovery for locked supervised iPhones
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "multiple_devices" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Recover passcodes for supervised iOS devices using data source
# ============================================================================
# Use case: Bulk recovery for all supervised iOS devices (with passcode escrow)
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_supervised_ios" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 4: Recover passcodes for supervised iPadOS devices
# ============================================================================
# Use case: Classroom iPads with forgotten passcodes
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_classroom_ipads" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 5: Recover passcode for specific user's supervised iOS device
# ============================================================================
# Use case: Executive's locked iPhone recovery
data "microsoft365_graph_beta_device_management_managed_device" "user_ios_device" {
  filter_type  = "user_id"
  filter_value = "user@example.com"
}

# Filter to only iOS supervised devices for this user
locals {
  supervised_ios_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.user_ios_device.items :
    device.id if device.operating_system == "iOS" && device.is_supervised == true
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_user_device" {

  device_ids = local.supervised_ios_devices

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 6: Recover passcode for device by serial number
# ============================================================================
# Use case: Device recovery using physical device serial number
data "microsoft365_graph_beta_device_management_managed_device" "device_by_serial" {
  filter_type  = "serial_number"
  filter_value = "DMQVG1234ABC"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_by_serial" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.device_by_serial.items : device.id]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# IMPORTANT NOTES
# ============================================================================
# 
# Passcode Escrow Requirements:
# - Passcodes MUST be escrowed during device enrollment
# - Supervised iOS/iPadOS devices typically escrow passcodes automatically
# - Check device enrollment profile settings for passcode escrow
# - Recovery will fail if passcode was never escrowed
# - If recovery fails, use reset_passcode action instead
#
# Platform Support:
# - iOS: Supported (supervised devices with escrow)
# - iPadOS: Supported (supervised devices with escrow)
# - macOS: Limited support (specific configurations only)
# - Windows: Not typically supported
# - Android: Not typically supported
#
# Recover vs Reset Passcode:
# - Recover: Retrieves existing escrowed passcode (no device change)
# - Reset: Generates new temporary passcode (device must unlock with new code)
# - Try recover first if device has passcode escrow
# - Use reset if recover fails or no escrow available
#
# Retrieving Recovered Passcode:
# 1. Navigate to Microsoft Intune admin center
# 2. Go to Devices > All devices
# 3. Select the device
# 4. View device properties
# 5. Look for passcode/recovery information
# 6. Securely communicate passcode to authorized user
#
# Security Considerations:
# - Recovered passcodes are sensitive credentials
# - Verify user identity before providing passcode
# - Communicate passcodes securely (not via email)
# - Document reason for passcode recovery
# - Follow organizational security policies
# - Ensure proper authorization before recovery
#
# Best Practices:
# - Only recover passcodes for authorized requests
# - Verify device ownership before recovery
# - Use with supervised devices for best results
# - Document business justification
# - Train help desk on passcode retrieval
# - Consider privacy implications
# - Monitor for repeated recovery requests
# - Implement approval workflow for recovery

