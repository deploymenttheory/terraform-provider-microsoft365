# Test 05: Get applicable content for a specific device

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get Entra ID devices
data "microsoft365_graph_beta_identity_and_access_device" "test_devices" {
  list_all = true
  timeouts = {
    read = "30s"
  }
}

# Create a deployment audience
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {}

# Wait for audience to propagate
resource "time_sleep" "wait_for_audience" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test]
  create_duration = "10s"
}

# Enroll devices in driver management
resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test" {
  update_category = "driver"
  entra_device_object_ids = [
    data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].id
  ]

  depends_on = [time_sleep.wait_for_audience]

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

# Add devices to the audience
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members" "test" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id
  member_type = "azureADDevice"

  members = [
    data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].id
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
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members.test]
  create_duration = "10s"
}

# Query applicable content - filter to specific device
data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id
  device_id   = data.microsoft365_graph_beta_identity_and_access_device.test_devices.items[0].id

  depends_on = [time_sleep.wait_for_members]

  timeouts = {
    read = "30s"
  }
}
