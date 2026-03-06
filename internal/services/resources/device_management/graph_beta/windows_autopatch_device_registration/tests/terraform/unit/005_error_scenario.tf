resource "microsoft365_graph_beta_device_management_windows_autopatch_device_registration" "test_005" {
  update_category = "feature"
  device_ids = [
    "invalid-device-id"
  ]
}
