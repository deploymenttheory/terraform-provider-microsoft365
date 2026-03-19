
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get managed devices from Intune which includes their Entra ID device IDs
data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  filter_type = "all"
  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_003" {
  update_category = "feature"
  # Use the azure_ad_device_id field which contains the Entra ID device object ID
  # Take up to 3 devices that have valid Entra ID device IDs
  entra_device_object_ids = [
    for device in slice(data.microsoft365_graph_beta_device_management_managed_device.test_devices.items, 0, min(3, length(data.microsoft365_graph_beta_device_management_managed_device.test_devices.items))) :
    device.azure_ad_device_id
    if device.azure_ad_device_id != null && device.azure_ad_device_id != ""
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
