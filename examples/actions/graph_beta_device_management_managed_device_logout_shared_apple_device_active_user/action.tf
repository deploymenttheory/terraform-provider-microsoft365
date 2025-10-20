# ============================================================================
# Example 1: Logout active user from single Shared iPad
# ============================================================================
# Use case: End of class period logout
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "single_device" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 2: Logout active users from multiple Shared iPads
# ============================================================================
# Use case: End of day logout for classroom cart devices
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "multiple_devices" {

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
# Example 3: Logout all Shared iPads in specific group
# ============================================================================
# Use case: Classroom management for scheduled logout
data "microsoft365_graph_beta_device_management_managed_device" "classroom_shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

# Filter to only devices with "SharediPad" in the name
locals {
  shared_ipad_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.classroom_shared_ipads.items :
    device.id if can(regex("SharediPad", device.device_name))
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_classroom" {

  device_ids = local.shared_ipad_devices

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 4: Logout Shared iPads by device name pattern
# ============================================================================
# Use case: Lab or cart devices with specific naming convention
data "microsoft365_graph_beta_device_management_managed_device" "lab_ipads" {
  filter_type  = "device_name"
  filter_value = "LAB-IPAD-"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_lab_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_ipads.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 5: Logout supervised iPads (potential Shared iPads)
# ============================================================================
# Use case: End of semester cleanup for all supervised iPads
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_supervised" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 6: Logout company-owned supervised iPads
# ============================================================================
# Use case: Institutional device rotation
data "microsoft365_graph_beta_device_management_managed_device" "company_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managedDeviceOwnerType eq 'company') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_company_ipads" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.company_ipads.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}
