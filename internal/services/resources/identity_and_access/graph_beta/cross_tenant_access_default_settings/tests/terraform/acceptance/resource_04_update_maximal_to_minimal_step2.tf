# ==============================================================================
# Cross-Tenant Access Default Settings - Update Maximal → Minimal (Step 2)
#
# Reduced to the minimal configuration. All optional blocks removed.
# User/group dependencies are no longer referenced by the cross-tenant resource
# and will be destroyed after this step.
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
