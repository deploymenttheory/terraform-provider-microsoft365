# Unit Test 02: List all enrolled devices

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all = true
}
