resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_group_sync_settings" "test" {
  tenant_id = "{{.TenantID}}"

  group_sync_inbound = {
    is_sync_allowed = true
  }
}
