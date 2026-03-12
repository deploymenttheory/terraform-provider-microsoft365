# ==============================================================================
# Cross-Tenant Access Default Settings - Minimal
#
# Configures the singleton cross-tenant access policy default settings with the
# smallest valid configuration: outbound B2B collaboration allowed for all users
# and all applications, with defaults restored on destroy.
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "test" {
  restore_defaults_on_destroy = true

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
