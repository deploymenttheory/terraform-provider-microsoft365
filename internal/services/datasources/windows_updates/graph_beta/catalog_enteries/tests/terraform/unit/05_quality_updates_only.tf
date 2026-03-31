data "microsoft365_graph_beta_windows_updates_catalog_enteries" "test" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}
