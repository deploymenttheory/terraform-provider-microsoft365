# Unit Test 07: Filter by driver update category

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all        = true
  update_category = "driver"
}
