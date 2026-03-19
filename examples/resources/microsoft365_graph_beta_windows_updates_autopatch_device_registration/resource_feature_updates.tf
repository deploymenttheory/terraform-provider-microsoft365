# Example: Register devices for Windows Autopatch feature updates
resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "feature_updates" {
  update_category = "feature"

  entra_device_object_ids = [
    "cccccccc-cccc-cccc-cccc-cccccccccccc",
    "dddddddd-dddd-dddd-dddd-dddddddddddd",
  ]
}
