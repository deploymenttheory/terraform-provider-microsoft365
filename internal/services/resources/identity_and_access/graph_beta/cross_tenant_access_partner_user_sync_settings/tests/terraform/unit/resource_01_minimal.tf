resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_user_sync_settings" "test" {
  tenant_id = "12345678-1234-1234-1234-123456789012"

  user_sync_inbound = {
    is_sync_allowed = true
  }
}
