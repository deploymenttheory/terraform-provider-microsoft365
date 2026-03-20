# Minimal example — only outbound B2B collaboration configured.
# All other settings remain at their partner defaults.
# On destroy, the partner configuration is soft deleted (can be restored within 30 days).

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = false

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
