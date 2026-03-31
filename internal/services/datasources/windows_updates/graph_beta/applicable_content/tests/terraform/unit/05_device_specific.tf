# Unit test: Get applicable content for specific device
data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id = "f660d844-30b7-46e4-a6cf-47e36164d3cb"
  device_id   = "fb95f07d-9e73-411d-99ab-7eca3a5122b1"
}
