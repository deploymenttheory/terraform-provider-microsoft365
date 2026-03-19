resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_004" {
  update_category = "feature"
  entra_device_object_ids = [
    "12345678-1234-1234-1234-123456789001",
    "12345678-1234-1234-1234-123456789002",
    "12345678-1234-1234-1234-123456789003"
  ]
}
