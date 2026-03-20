# ==============================================================================
# Cross-Tenant Access Partner Settings - Update Minimal to Maximal (Step 1)
#
# Step 1: Deploy minimal configuration with only outbound B2B collaboration.
# This will be expanded to the full maximal configuration in step 2.
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
