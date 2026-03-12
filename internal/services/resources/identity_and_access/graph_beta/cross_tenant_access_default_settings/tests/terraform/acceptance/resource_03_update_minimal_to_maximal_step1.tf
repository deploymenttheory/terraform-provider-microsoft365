# ==============================================================================
# Cross-Tenant Access Default Settings - Update Minimal → Maximal (Step 1)
#
# Initial minimal configuration. Step 2 will expand this to the full maximal
# configuration, exercising the PATCH update path.
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
