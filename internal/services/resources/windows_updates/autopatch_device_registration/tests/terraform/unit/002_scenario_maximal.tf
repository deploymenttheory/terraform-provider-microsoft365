resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_002" {
  update_category = "quality"
  device_ids = [
    "12345678-1234-1234-1234-123456789001",
    "12345678-1234-1234-1234-123456789002",
    "12345678-1234-1234-1234-123456789003"
  ]
}
