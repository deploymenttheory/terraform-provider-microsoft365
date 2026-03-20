# ==============================================================================
# Cross-Tenant Access Partner Settings - Update Maximal to Minimal (Step 1)
#
# Step 1: Deploy maximal configuration with all blocks populated.
# This will be reduced to the minimal configuration in step 2.
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "test" {
  tenant_id   = "a22ff489-2ea9-48de-8d58-fa130b532d5d"
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

  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }

  tenant_restrictions = {
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

  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
