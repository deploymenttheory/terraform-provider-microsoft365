data "microsoft365_graph_beta_windows_updates_catalog_enteries" "test" {
  filter_type  = "display_name"
  filter_value = "SecurityUpdate"
}
