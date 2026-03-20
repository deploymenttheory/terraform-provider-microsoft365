# B2B collaboration — controls inbound and outbound B2B guest access for a specific partner.
# Both directions are configured to allow all users and all applications.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

  b2b_collaboration_inbound = {
    users_and_groups = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }

  b2b_collaboration_outbound = {
    users_and_groups = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }
}
