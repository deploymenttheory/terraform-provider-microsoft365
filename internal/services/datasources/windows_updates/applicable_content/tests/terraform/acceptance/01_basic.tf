# Test 01: Get applicable content for a deployment audience
# This test creates a full dependency chain to test applicable content retrieval

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get managed devices from Intune
data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
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
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_ad_device_id
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
    data.microsoft365_graph_beta_device_management_managed_device.test_devices.items[0].azure_ad_device_id
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

# Query applicable content
data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id

  depends_on = [time_sleep.wait_for_members]

  timeouts = {
    read = "30s"
  }
}

output "applicable_content_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.test.applicable_content)
}
