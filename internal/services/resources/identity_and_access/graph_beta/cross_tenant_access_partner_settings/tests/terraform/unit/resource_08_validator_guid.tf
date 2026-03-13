resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "test" {
  tenant_id = "12345678-1234-1234-1234-123456789012"

  b2b_direct_connect_outbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = "11111111-1111-1111-1111-111111111111"
          target_type = "user"
        },
        {
          target      = "33333333-3333-3333-3333-333333333333"
          target_type = "group"
        }
      ]
    }
    applications = {
      access_type = "blocked"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }
}
