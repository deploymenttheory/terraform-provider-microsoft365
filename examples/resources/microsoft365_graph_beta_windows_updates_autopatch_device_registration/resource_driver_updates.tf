# Example: Register devices for Windows Autopatch driver updates
resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "driver_updates" {
  update_category = "driver"

  entra_device_object_ids = [
    "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee",
  ]
}
