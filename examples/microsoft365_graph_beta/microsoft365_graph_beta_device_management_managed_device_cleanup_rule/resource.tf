# Example 1: Android Device Administrator cleanup rule after 90 days of inactivity
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "android_da_cleanup" {
  display_name                                = "Android Device Admin Cleanup - 90 Days"
  description                                 = "Automatically remove Android Device Administrator devices that haven't contacted Intune for 90 days"
  device_cleanup_rule_platform_type           = "androidDeviceAdministrator"
  device_inactivity_before_retirement_in_days = 90

  timeouts {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 2: iOS device cleanup rule after 60 days
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "ios_cleanup" {
  display_name                                = "iOS Device Cleanup - 60 Days"
  description                                 = "Remove iOS devices inactive for 60 days"
  device_cleanup_rule_platform_type           = "ios"
  device_inactivity_before_retirement_in_days = 60
}

# Example 3: Windows device cleanup rule after 180 days
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "windows_cleanup" {
  display_name                                = "Windows Device Cleanup - 180 Days"
  description                                 = "Remove Windows devices that haven't synced for 6 months"
  device_cleanup_rule_platform_type           = "windows"
  device_inactivity_before_retirement_in_days = 180
}

# Example 4: All platforms cleanup rule after 365 days
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "all_platforms_cleanup" {
  display_name                                = "All Devices Cleanup - 1 Year"
  description                                 = "Annual cleanup for all device types"
  device_cleanup_rule_platform_type           = "all"
  device_inactivity_before_retirement_in_days = 365
}

# Example 5: ChromeOS device cleanup rule with minimal inactivity period
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "chromeos_immediate" {
  display_name                                = "ChromeOS Immediate Cleanup"
  device_cleanup_rule_platform_type           = "chromeOS"
  device_inactivity_before_retirement_in_days = 0
}

# Example 6: Android work profile cleanup rule
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "android_work_profile" {
  display_name                                = "Android Work Profile Cleanup"
  description                                 = "Cleanup for personally owned devices with work profiles"
  device_cleanup_rule_platform_type           = "androidPersonallyOwnedWorkProfile"
  device_inactivity_before_retirement_in_days = 120
}

# Example 7: Windows Holographic cleanup rule
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "hololens_cleanup" {
  display_name                                = "HoloLens Device Cleanup"
  description                                 = "Remove HoloLens devices after extended inactivity"
  device_cleanup_rule_platform_type           = "windowsHolographic"
  device_inactivity_before_retirement_in_days = 30
}