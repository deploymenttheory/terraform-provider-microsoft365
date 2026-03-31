# Acceptance Test 02: List all enrolled devices
# This test registers multiple devices and then lists all enrolled devices

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Get Test Devices
# ==============================================================================

data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  odata_query = "$filter=operatingSystem eq 'Windows'&$top=3"
}

# ==============================================================================
# Register Devices for Windows Updates
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test" {
  update_category = "quality"
  entra_device_object_ids = [
    for device in slice(data.microsoft365_graph_beta_device_management_managed_device.test_devices.items, 0, min(3, length(data.microsoft365_graph_beta_device_management_managed_device.test_devices.items))) :
    device.azure_ad_device_id
    if device.azure_ad_device_id != null && device.azure_ad_device_id != ""
  ]
}

# ==============================================================================
# Wait for Device Registration
# ==============================================================================

resource "time_sleep" "wait_for_registration" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_device_registration.test]
  create_duration = "10s"
}

# ==============================================================================
# Data Source - List All Devices
# ==============================================================================

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all = true

  depends_on = [time_sleep.wait_for_registration]
}

# ==============================================================================
# Outputs
# ==============================================================================

output "total_devices" {
  value = length(data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices)
}

output "device_ids" {
  value = [for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices : device.id]
}
