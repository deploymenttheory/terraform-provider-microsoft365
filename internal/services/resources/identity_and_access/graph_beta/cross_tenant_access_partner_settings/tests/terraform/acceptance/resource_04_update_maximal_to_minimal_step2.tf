# ==============================================================================
# Cross-Tenant Access Partner Settings - Update Maximal to Minimal (Step 2)
#
# Step 2: Reduce to minimal configuration with only outbound B2B collaboration.
# This demonstrates the PATCH update capability when removing configuration blocks.
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
