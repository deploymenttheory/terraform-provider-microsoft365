resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = "minimal-catalog-entry-id"
    catalog_entry_type = "featureUpdate"
  }
}
