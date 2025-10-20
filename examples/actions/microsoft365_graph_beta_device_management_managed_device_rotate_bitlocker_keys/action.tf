# Data source to find all Windows devices with BitLocker enabled
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter = "operatingSystem eq 'Windows'"
}

# Example 1: Rotate BitLocker keys on specific Windows managed devices
# Use this for targeted key rotation on specific devices (e.g., after security incident)
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_specific_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]
}

# Example 2: Rotate BitLocker keys on co-managed Windows devices
# Use this for devices managed by both Intune and Configuration Manager
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_comanaged_devices" {
  comanaged_device_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

# Example 3: Rotate BitLocker keys on both managed and co-managed devices
# Use this for mixed device management scenarios
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_mixed_devices" {
  managed_device_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  ]

  comanaged_device_ids = [
    "cccccccc-cccc-cccc-cccc-cccccccccccc"
  ]
}

# Example 4: Rotate BitLocker keys on all Windows 10 devices
# Use this for scheduled maintenance or compliance requirement across all Windows 10 devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_10_devices" {
  filter = "operatingSystem eq 'Windows' and contains(osVersion, '10.')"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_all_windows_10" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_10_devices.managed_devices : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 5: Rotate BitLocker keys on Windows 11 Enterprise devices
# Use this for targeting specific Windows editions
data "microsoft365_graph_beta_device_management_managed_device" "windows_11_enterprise" {
  filter = "operatingSystem eq 'Windows' and contains(osVersion, '11.') and skuFamily eq 'Windows.Enterprise'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_windows_11_enterprise" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_11_enterprise.managed_devices : device.id]
}

# Example 6: Rotate BitLocker keys on non-compliant Windows devices
# Use this as part of compliance remediation process
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_windows" {
  filter = "complianceState eq 'noncompliant' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_noncompliant" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_windows.managed_devices : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Rotate BitLocker keys on devices by name pattern
# Use this when you need to target specific departments or device groups
data "microsoft365_graph_beta_device_management_managed_device" "finance_dept_devices" {
  filter = "startswith(deviceName, 'FIN-') and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_finance_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_dept_devices.managed_devices : device.id]
}

# Example 8: Rotate BitLocker keys on recently enrolled devices
# Use this to ensure new devices have keys properly escrowed after initial setup
data "microsoft365_graph_beta_device_management_managed_device" "recently_enrolled_windows" {
  filter = "enrolledDateTime gt 2024-01-01T00:00:00Z and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_new_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.recently_enrolled_windows.managed_devices : device.id]
}

# Example 9: Rotate BitLocker keys on Azure AD joined Windows devices
# Use this to target cloud-native Windows devices specifically
data "microsoft365_graph_beta_device_management_managed_device" "azure_ad_joined_windows" {
  filter = "joinType eq 'azureADJoined' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_azure_ad_joined" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.azure_ad_joined_windows.managed_devices : device.id]
}

# Example 10: Rotate BitLocker keys on corporate-owned Windows devices
# Use this to differentiate between corporate and BYOD devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter = "managedDeviceOwnerType eq 'company' and operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_corporate_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.managed_devices : device.id]
}

# Example 11: Rotate BitLocker keys with custom extended timeout
# Use this for very large-scale operations (500+ devices)
action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "rotate_with_extended_timeout" {
  managed_device_ids = [
    "device-id-1",
    "device-id-2",
    "device-id-3"
  ]

  timeouts = {
    invoke = "45m"
  }
}

# Example 12: Scheduled BitLocker key rotation for compliance
# Use this for regular security maintenance (e.g., quarterly key rotation)
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_managed" {
  filter = "operatingSystem eq 'Windows' and managementAgent eq 'mdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_rotate_bitlocker_keys" "quarterly_key_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_managed.managed_devices : device.id]

  timeouts = {
    invoke = "60m"
  }
}

