# Unit Test 08: Look up device enrollment by device name

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  device_name = "TEST-DEVICE-001"
}
