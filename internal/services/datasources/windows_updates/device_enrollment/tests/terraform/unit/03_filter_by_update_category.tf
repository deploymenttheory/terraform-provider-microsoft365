# Unit Test 03: List all devices filtered by update category

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all        = true
  update_category = "quality"
}
