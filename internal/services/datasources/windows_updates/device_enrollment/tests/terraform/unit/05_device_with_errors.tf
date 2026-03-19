# Unit Test 05: Device with registration errors

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  entra_device_id = "0ee3eb63-caf3-44ce-9769-b83188cc683d"
}
