# B2B direct connect — blocks direct connect for both inbound and outbound directions.
# This prevents Teams Connect shared channels with the specified partner tenant.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

  b2b_direct_connect_inbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
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

  b2b_direct_connect_outbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
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
