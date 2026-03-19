# Unit Test 06: Filter by feature update category

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all        = true
  update_category = "feature"
}
