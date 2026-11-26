# Test 04: Get only content with matched devices

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get Entra ID devices
data "microsoft365_graph_beta_identity_and_access_device" "test_devices" {
  odata_query = "operatingSystem eq 'Windows' and accountEnabled eq true"
  timeouts = {
    read = "30s"
  }
}

# Enroll devices in driver management
resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test" {
  update_category = "driver"
  entra_device_object_ids = [
    data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].device_id
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Wait for device registration
resource "time_sleep" "wait_for_registration" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_device_registration.test]
  create_duration = "10s"
}

# Create an audience with the registered device
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {
  member_type = "azureADDevice"

  members = [
    data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].device_id
  ]

  depends_on = [time_sleep.wait_for_registration]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Wait for audience members to propagate
resource "time_sleep" "wait_for_members" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test]
  create_duration = "10s"
}

# Query applicable content - only content with matched devices
data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id        = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id
  include_no_matches = false

  depends_on = [time_sleep.wait_for_members]

  timeouts = {
    read = "30s"
  }
}
