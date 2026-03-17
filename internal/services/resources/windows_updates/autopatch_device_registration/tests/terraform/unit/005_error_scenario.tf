resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_005" {
  update_category = "feature"
  device_ids = [
    "invalid-device-id"
  ]
}
