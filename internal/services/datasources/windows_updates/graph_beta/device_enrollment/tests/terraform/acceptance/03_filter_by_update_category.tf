# Acceptance Test 03: Filter devices by update category
# This test lists all devices filtered by a specific update category

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
# Data Source - Filter by Quality Update Category
# ==============================================================================

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all        = true
  update_category = "quality"

  depends_on = [time_sleep.wait_for_registration]
}

# ==============================================================================
# Outputs
# ==============================================================================

output "filtered_devices" {
  value = length(data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices)
}

output "enrollments" {
  value = [for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.test.devices : device.enrollments]
}
