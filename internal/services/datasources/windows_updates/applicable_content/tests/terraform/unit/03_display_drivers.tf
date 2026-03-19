# Unit test: Get only display driver updates
data "microsoft365_graph_beta_windows_updates_applicable_content" "test" {
  audience_id        = "f660d844-30b7-46e4-a6cf-47e36164d3cb"
  catalog_entry_type = "driver"
  driver_class       = "Display"
}
