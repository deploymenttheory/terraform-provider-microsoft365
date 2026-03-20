# ==============================================================================
# Cross-Tenant Access Partner Settings - Minimal
#
# Configures partner-specific cross-tenant access settings with the smallest
# valid configuration: outbound B2B collaboration allowed for all users and
# all applications, with hard delete enabled on destroy.
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "test" {
  tenant_id   = "a22ff489-2ea9-48de-8d58-fa130b532d5d"
  hard_delete = true

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
