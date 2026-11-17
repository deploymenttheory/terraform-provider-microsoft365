# Autopatch membership uses Entra deviceId values and does not require Intune enrollment.
data "microsoft365_graph_beta_identity_and_access_device" "devices" {
  odata_query = "operatingSystem eq 'Windows' and accountEnabled eq true"

  timeouts = {
    read = "30s"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "example" {
  entra_device_object_ids = data.microsoft365_graph_beta_identity_and_access_device.devices.items[*].device_id

  lifecycle {
    precondition {
      condition     = length(data.microsoft365_graph_beta_identity_and_access_device.devices.items) > 0
      error_message = "At least one enabled Windows Entra device is required. Intune enrollment is not required."
    }
  }

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
