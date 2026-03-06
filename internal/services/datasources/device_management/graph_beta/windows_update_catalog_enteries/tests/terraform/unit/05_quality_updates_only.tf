data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "test" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}
