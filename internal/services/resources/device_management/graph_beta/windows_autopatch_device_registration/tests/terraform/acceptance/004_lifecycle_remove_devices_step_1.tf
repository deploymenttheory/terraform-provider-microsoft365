
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Get managed devices from Intune which includes their Entra ID device IDs
data "microsoft365_graph_beta_device_management_managed_device" "test_devices" {
  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_autopatch_device_registration" "test_004" {
  update_category = "feature"
  # Use the azure_ad_device_id field which contains the Entra ID device object ID
  # Take up to 3 devices that have valid Entra ID device IDs
  device_ids = [
    for device in slice(data.microsoft365_graph_beta_device_management_managed_device.test_devices.managed_devices, 0, min(3, length(data.microsoft365_graph_beta_device_management_managed_device.test_devices.managed_devices))) :
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
