resource "microsoft365_graph_beta_device_management_windows_autopatch_device_registration" "test_004" {
  update_category = "feature"
  device_ids = [
    "12345678-1234-1234-1234-123456789001",
    "12345678-1234-1234-1234-123456789002",
    "12345678-1234-1234-1234-123456789003"
  ]
}
