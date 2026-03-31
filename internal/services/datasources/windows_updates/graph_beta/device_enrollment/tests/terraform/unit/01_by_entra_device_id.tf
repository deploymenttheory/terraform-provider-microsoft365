# Unit Test 01: Look up device enrollment by Entra device ID

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  entra_device_id = "fb95f07d-9e73-411d-99ab-7eca3a5122b1"
}
