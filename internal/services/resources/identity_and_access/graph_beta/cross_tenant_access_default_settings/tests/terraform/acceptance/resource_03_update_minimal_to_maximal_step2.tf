# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# User Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "acc_test_user_1" {
  display_name        = "acc-test-cta-user1-${random_string.suffix.result}"
  user_principal_name = "acc-test-cta-user1-${random_string.suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-cta-user1-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-cta-group1-${random_string.suffix.result}"
  mail_nickname    = "acc-test-cta-group1-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for cross-tenant access default settings outbound direct connect"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-cta-group2-${random_string.suffix.result}"
  mail_nickname    = "acc-test-cta-group2-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for cross-tenant access default settings outbound direct connect"
  hard_delete      = true
}

# ==============================================================================
# Time Sleep
# ==============================================================================

resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    microsoft365_graph_beta_users_user.acc_test_user_1,
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2,
  ]
  create_duration = "30s"
}

# ==============================================================================
# Cross-Tenant Access Default Settings - Update Minimal → Maximal (Step 2)
#
# All blocks are now configured. The PATCH update path is exercised against
# the existing state from Step 1.
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "test" {
  restore_defaults_on_destroy = true

  depends_on = [time_sleep.wait_30_seconds]

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
          target      = "Office365"
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
          target      = microsoft365_graph_beta_users_user.acc_test_user_1.id
          target_type = "user"
        },
        {
          target      = microsoft365_graph_beta_groups_group.acc_test_group_1.id
          target_type = "group"
        },
        {
          target      = microsoft365_graph_beta_groups_group.acc_test_group_2.id
          target_type = "group"
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

  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
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
