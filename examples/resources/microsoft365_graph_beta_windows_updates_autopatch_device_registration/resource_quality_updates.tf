# Example: Register devices for Windows Autopatch quality updates
resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "quality_updates" {
  update_category = "quality"

  entra_device_object_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
  ]
}
