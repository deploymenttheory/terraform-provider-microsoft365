# A configured empty set clears the list in the tenant.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  excluded_file_extensions_for_sync_app = []
}
