# REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta

# Data source to find iOS devices with VPP apps assigned
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices_with_vpp" {
  filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

# Example 1: Revoke Apple VPP licenses from specific iOS managed devices
# Use this when reclaiming VPP licenses from specific devices that are being retired
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_specific_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]
}

# Example 2: Revoke Apple VPP licenses from co-managed iOS devices
# Use this for devices managed by both Intune and Configuration Manager
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_comanaged_devices" {
  comanaged_device_ids = [
    "11111111-1111-1111-1111-111111111111"
  ]
}

# Example 3: Revoke Apple VPP licenses from both managed and co-managed devices
# Use this for mixed device management scenarios
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_mixed_devices" {
  managed_device_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  ]

  comanaged_device_ids = [
    "cccccccc-cccc-cccc-cccc-cccccccccccc"
  ]
}

# Example 4: Revoke Apple VPP licenses from all iOS devices using data source
# Use this when performing bulk license reclamation for all iOS/iPadOS devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_all_ios_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices_with_vpp.managed_devices : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Revoke licenses from specific devices by device name using filter
# Use this when you need to target devices by name for license recovery
data "microsoft365_graph_beta_device_management_managed_device" "specific_ios_devices" {
  filter = "startswith(deviceName, 'iPad-Retail-') and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_by_name_pattern" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.specific_ios_devices.managed_devices : device.id]
}

# Example 6: Revoke licenses from non-compliant iOS devices
# Use this to reclaim licenses from devices that are no longer compliant
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_ios" {
  filter = "complianceState eq 'noncompliant' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_noncompliant" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_ios.managed_devices : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 7: Revoke licenses from inactive iOS devices (not synced in 30+ days)
# Use this to optimize license allocation by reclaiming from inactive devices
data "microsoft365_graph_beta_device_management_managed_device" "inactive_ios_devices" {
  filter = "lastSyncDateTime lt 2024-01-01T00:00:00Z and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_inactive" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.inactive_ios_devices.managed_devices : device.id]
}

# Example 8: Revoke licenses from supervised iOS devices
# Use this when you need to reclaim licenses specifically from supervised devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter = "isSupervised eq true and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_supervised" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.managed_devices : device.id]
}

# Example 9: Revoke licenses with custom timeout
# Use this for large-scale license revocation operations
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_with_timeout" {
  managed_device_ids = [
    "device-id-1",
    "device-id-2",
    "device-id-3"
  ]

  timeouts = {
    invoke = "30m"
  }
}

