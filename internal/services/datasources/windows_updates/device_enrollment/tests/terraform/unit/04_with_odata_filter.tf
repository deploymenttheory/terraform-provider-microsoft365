# Unit Test 04: List devices with custom OData filter

data "microsoft365_graph_beta_windows_updates_device_enrollment" "test" {
  list_all     = true
  odata_filter = "id eq 'fb95f07d-9e73-411d-99ab-7eca3a5122b1'"
}
