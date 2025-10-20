# Example 1: Bypass Activation Lock for a single device
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_single_device" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Bypass Activation Lock for multiple devices (batch processing)
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Bypass Activation Lock for supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS') and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_supervised_ios" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Bypass Activation Lock for DEP-enrolled macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "dep_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (deviceEnrollmentType eq 'deviceEnrollmentProgram')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_dep_macos" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dep_macos.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Bypass Activation Lock for departing employee's Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_apple_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_departing_user" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_apple_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Bypass Activation Lock for corporate-owned Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_apple" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS') and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_corporate_apple" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_apple.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Bypass Activation Lock for devices with specific model (e.g., iPhone 13)
data "microsoft365_graph_beta_device_management_managed_device" "iphone_13_devices" {
  filter_type  = "odata"
  odata_filter = "model eq 'iPhone 13' and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_iphone_13" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.iphone_13_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "bypassed_device_count" {
  value       = length(action.bypass_batch.device_ids)
  description = "Number of devices that received Activation Lock bypass command"
}

output "bypassed_corporate_count" {
  value       = length(action.bypass_corporate_apple.device_ids)
  description = "Number of corporate Apple devices with bypass codes generated"
}

# Important Notes:
# 
# What is Activation Lock?
# - Apple security feature that prevents device reactivation after factory reset
# - Automatically enabled when Find My iPhone/iPad/Mac is turned on
# - Requires original Apple ID and password to reactivate device
# - Protects against unauthorized device use after theft or loss
#
# When to Use Bypass Activation Lock:
# - Preparing corporate devices for reassignment to new employees
# - Recovering devices from departing employees who forgot to disable Find My
# - Preparing devices for return to vendor or recycling
# - Handling devices with lost or forgotten Apple ID credentials
# - Bulk device refresh/replacement projects
# - Emergency device recovery scenarios
#
# Platform Requirements:
# - iOS/iPadOS: Devices MUST be supervised (enrolled via DEP/ABM or Apple Configurator)
# - macOS: Devices should be enrolled via DEP/ABM (Automated Device Enrollment)
# - Windows/Android: Not supported (Activation Lock is Apple-only feature)
#
# How Activation Lock Bypass Works:
# 1. Issue bypass command via this action
# 2. Intune generates a unique bypass code for each device
# 3. Bypass code is stored in Intune device properties
# 4. Retrieve bypass code from Intune admin portal
# 5. Factory reset/wipe the device
# 6. During setup, device shows Activation Lock screen
# 7. Enter bypass code in password field to bypass lock
# 8. Device completes setup without requiring user's Apple ID
#
# Workflow Example:
# Step 1: Employee leaves organization, device in their possession
# Step 2: IT issues bypass activation lock command (this action)
# Step 3: IT retrieves bypass code from Intune portal
# Step 4: IT wipes device remotely or obtains physical possession
# Step 5: During device setup, Activation Lock screen appears
# Step 6: IT enters bypass code to unlock device
# Step 7: Device can now be re-enrolled and assigned to new user
#
# Retrieving Bypass Codes:
# - Navigate to Intune admin center (https://intune.microsoft.com)
# - Go to Devices > All devices
# - Select the device
# - Under "Hardware" section, find "Activation Lock bypass code"
# - Copy code and securely store or use immediately
# - Code format: Usually 6-8 alphanumeric characters
#
# Security Considerations:
# - Bypass codes are sensitive credentials - treat like passwords
# - Only authorized IT staff should have access to codes
# - Document bypass code usage for compliance/audit purposes
# - Consider implementing approval workflow for bypass requests
# - Verify device ownership before issuing bypass
# - Bypass codes should be stored in secure password vault
# - Rotate or regenerate codes periodically if stored long-term
#
# Limitations:
# - Cannot bypass Activation Lock on personal (non-supervised) iOS/iPadOS devices
# - Cannot bypass on macOS devices not enrolled via DEP/ABM (limited support)
# - Bypass code only works for specific device it was generated for
# - Must retrieve bypass code before device is erased (code stored in Intune)
# - Some older device models may not support Activation Lock bypass
# - Device must have had Activation Lock enabled for bypass to be relevant
#
# Best Practices:
# - Issue bypass command BEFORE wiping device when possible
# - Store bypass codes securely in enterprise password manager
# - Document which devices have bypass codes generated
# - Include Activation Lock bypass in offboarding procedures
# - Test bypass process in controlled environment first
# - Verify device supervision status before attempting bypass
# - Consider enabling automatic bypass code escrow during enrollment
# - Train help desk staff on bypass code retrieval and usage
# - Maintain audit log of bypass code usage
#
# Common Issues and Solutions:
# 
# Issue: Bypass command fails with "Device not supervised"
# Solution: iOS/iPadOS devices must be supervised. Re-enroll via DEP/ABM or Apple Configurator
#
# Issue: Bypass command fails with "Activation Lock not enabled"
# Solution: Device doesn't have Find My enabled. No bypass needed for this device
#
# Issue: Bypass code doesn't work during device setup
# Solution: Verify you copied the code correctly. Try entering code in all caps or lowercase
#
# Issue: Can't find bypass code in Intune portal
# Solution: Code may take a few minutes to appear. Refresh device properties page
#
# Issue: macOS device doesn't accept bypass code
# Solution: Verify device was enrolled via DEP/ABM. Manually enrolled Macs have limited support
#
# Integration with Other Actions:
# - Often used before or after wipe action
# - Can be combined with retire for less aggressive device cleanup
# - May be used with disable lost mode if device is in lost mode
# - Should be part of comprehensive device lifecycle management
#
# Compliance and Legal Considerations:
# - Ensure you have legal right to bypass device (corporate ownership)
# - Document business justification for bypass in audit logs
# - Consider privacy implications in different jurisdictions
# - Review employment agreements regarding device management
# - Maintain records of device ownership and bypass authorization
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta

