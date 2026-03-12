resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_user_sync_settings" "test" {
  tenant_id    = "{{.TenantID}}"
  display_name = "Partner Sync Configuration"

  user_sync_inbound = {
    is_sync_allowed = true
  }
}
