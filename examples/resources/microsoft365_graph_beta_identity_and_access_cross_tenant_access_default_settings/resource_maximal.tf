# Maximal example — all available blocks configured.
#
# b2b_direct_connect_outbound references specific users and groups created as
# dependencies. A 30-second time_sleep ensures the objects have fully propagated
# in the directory before the policy is patched with their IDs.

resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "example-cta-user"
  user_principal_name = "example-cta-user@contoso.com"
  mail_nickname       = "example-cta-user"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_groups_group" "example_group_1" {
  display_name     = "example-cta-group-1"
  mail_nickname    = "example-cta-group-1"
  mail_enabled     = false
  security_enabled = true
  description      = "Example group 1 for outbound direct connect blocking"
}

resource "microsoft365_graph_beta_groups_group" "example_group_2" {
  display_name     = "example-cta-group-2"
  mail_nickname    = "example-cta-group-2"
  mail_enabled     = false
  security_enabled = true
  description      = "Example group 2 for outbound direct connect blocking"
}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    microsoft365_graph_beta_users_user.example_user,
    microsoft365_graph_beta_groups_group.example_group_1,
    microsoft365_graph_beta_groups_group.example_group_2,
  ]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
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

  # b2b_direct_connect_inbound only supports "AllUsers" for users_and_groups.
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

  # b2b_direct_connect_outbound supports specific user and group GUIDs.
  b2b_direct_connect_outbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = microsoft365_graph_beta_users_user.example_user.id
          target_type = "user"
        },
        {
          target      = microsoft365_graph_beta_groups_group.example_group_1.id
          target_type = "group"
        },
        {
          target      = microsoft365_graph_beta_groups_group.example_group_2.id
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
