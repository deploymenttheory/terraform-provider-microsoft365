# ============================================================================
# Example 1: Clean single Windows device (remove user data)
# ============================================================================
# Use case: Full cleanup before reassignment
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_full_clean" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]

  keep_user_data = false

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 2: Clean single Windows device (preserve user data)
# ============================================================================
# Use case: Application cleanup while keeping user profiles
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_preserve_data" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]

  keep_user_data = true

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Clean multiple Windows devices (remove user data)
# ============================================================================
# Use case: Bulk device refresh before new deployment
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "multiple_devices_full_clean" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  keep_user_data = false

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 4: Clean non-compliant Windows devices (remove user data)
# ============================================================================
# Use case: Remediate non-compliant devices with full cleanup
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_noncompliant" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : device.id]

  keep_user_data = false

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 5: Clean Windows 11 devices (preserve user data)
# ============================================================================
# Use case: Application troubleshooting while preserving user environment
data "microsoft365_graph_beta_device_management_managed_device" "windows11_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (startswith(osVersion, '10.0.22'))"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_win11_preserve_data" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows11_devices.items : device.id]

  keep_user_data = true

  timeouts = {
    invoke = "25m"
  }
}

# ============================================================================
# Example 6: Clean corporate Windows devices (remove user data)
# ============================================================================
# Use case: Company-owned device refresh
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_corporate_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

  keep_user_data = false

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# IMPORTANT NOTES
# ============================================================================
# 
# Clean vs Wipe vs Retire:
# - Clean: Removes apps/settings, optionally keeps user data, device stays enrolled
# - Wipe: Factory reset, removes all data, device must re-enroll
# - Retire: Removes company data only, preserves personal data
#
# keep_user_data Parameter (REQUIRED):
# - false: Removes user profiles and data (full cleanup)
# - true: Preserves user profiles and data (app cleanup only)
# - Must be explicitly set - no default value
#
# Platform Support:
# - Windows 10: Fully supported
# - Windows 11: Fully supported
# - Other platforms: Not supported (Windows-only action)
#
# What Gets Removed (when keep_user_data = false):
# - Installed applications (except inbox Windows apps)
# - User profiles
# - User data (documents, desktop, etc.)
# - Device configuration settings
# - Company policies and profiles
#
# What Gets Removed (when keep_user_data = true):
# - Installed applications (except inbox Windows apps)
# - Device configuration settings
# - Company policies and profiles
# - User profiles and data are PRESERVED
#
# What is Preserved (Always):
# - Windows OS installation
# - Intune enrollment
# - Device in Intune management
# - Inbox Windows applications
# - Device hardware configuration
#
# Common Use Cases:
# - keep_user_data = false:
#   * Device reassignment to new user
#   * Malware removal requiring full cleanup
#   * Compliance remediation
#   * Device refresh before new deployment
#   * Preparing device for return/disposal
#
# - keep_user_data = true:
#   * Application troubleshooting
#   * Software bloat removal
#   * Configuration reset
#   * Maintaining user environment
#   * Performance optimization
#
# Best Practices:
# - Test with small device groups first
# - Schedule during maintenance windows
# - Notify users before cleaning their devices
# - Use keep_user_data=true for user devices when possible
# - Use keep_user_data=false for device reassignment
# - Document business justification
# - Monitor device status after clean
# - Allow sufficient timeout (devices may take 10-20 minutes)
#
# User Impact:
# - Device will be unavailable during clean process
# - Unsaved work will be lost
# - Process typically takes 10-20 minutes
# - Device remains enrolled (no re-enrollment needed)
# - User can log back in after clean completes
# - With keep_user_data=true, user profile is intact
# - With keep_user_data=false, user profile is removed
#
# Prerequisites:
# - Device must be Windows 10 or Windows 11
# - Device must be online
# - Device must be enrolled in Intune
# - Sufficient admin permissions required
#
# After Clean Operation:
# - Device remains in Intune
# - Policies will reapply automatically
# - Applications must be reinstalled
# - Check device status in Intune admin center
# - User may need to reconfigure personal settings
# - Company Portal remains available

