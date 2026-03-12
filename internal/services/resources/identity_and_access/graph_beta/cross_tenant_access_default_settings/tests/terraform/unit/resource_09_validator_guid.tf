resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "test" {
  restore_defaults_on_destroy = true

  b2b_direct_connect_outbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = "11111111-1111-1111-1111-111111111111"
          target_type = "user"
        },
        {
          target      = "22222222-2222-2222-2222-222222222222"
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
