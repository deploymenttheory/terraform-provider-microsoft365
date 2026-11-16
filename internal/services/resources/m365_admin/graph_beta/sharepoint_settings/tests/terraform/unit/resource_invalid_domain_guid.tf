# Invalid: allowed_domain_guids_for_sync_app values must be lower-case GUIDs.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  allowed_domain_guids_for_sync_app = ["BDD1AB9B-3FD0-4DEF-A761-EC8D7471732C"]
}
