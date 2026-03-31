# Acceptance Test 01: Look up device enrollment by Entra device ID
# This test registers a device and then looks it up by its Entra device ID

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
  odata_query = "$filter=operatingSystem eq 'Windows'&$top=1"
}

# ==============================================================================
# Register Device for Windows Updates
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test" {
  update_category = "quality"
  entra_device_object_ids = [
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_ad_device_id
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
# Data Source - Lookup by Entra Device ID
# ==============================================================================

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  entra_device_id = data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_ad_device_id

  depends_on = [time_sleep.wait_for_registration]
}

# ==============================================================================
# Outputs
# ==============================================================================

output "device_id" {
  value = data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices[0].id
}

output "enrollments" {
  value = data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices[0].enrollments
}
