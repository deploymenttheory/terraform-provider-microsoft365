# Complete example — all available blocks and settings configured together.
#
# This resource manages the singleton default cross-tenant access policy
# configuration (crossTenantAccessPolicyConfigurationDefault) that applies
# to all external tenants unless overridden by a partner-specific policy.
#
# Because no POST/DELETE endpoints exist, the provider issues a PATCH on create
# and update, and optionally calls resetToSystemDefault on destroy when
# restore_defaults_on_destroy = true.
#
# Dependencies:
#   b2b_direct_connect_outbound.users_and_groups.targets references specific
#   user and group IDs. A time_sleep ensures those objects are fully propagated
#   in Entra ID before the policy PATCH is issued.

# ==============================================================================
# User Dependency
# ==============================================================================

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

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "example_group_1" {
  display_name     = "example-cta-group-1"
  mail_nickname    = "example-cta-group-1"
  mail_enabled     = false
  security_enabled = true
  description      = "Group 1 blocked from outbound direct connect"
}

resource "microsoft365_graph_beta_groups_group" "example_group_2" {
  display_name     = "example-cta-group-2"
  mail_nickname    = "example-cta-group-2"
  mail_enabled     = false
  security_enabled = true
  description      = "Group 2 blocked from outbound direct connect"
}

# ==============================================================================
# Propagation Wait
# ==============================================================================

resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    microsoft365_graph_beta_users_user.example_user,
    microsoft365_graph_beta_groups_group.example_group_1,
    microsoft365_graph_beta_groups_group.example_group_2,
  ]
  create_duration = "30s"
}

# ==============================================================================
# Cross-Tenant Access Default Settings
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  # When true, destroys the resource by calling resetToSystemDefault and
  # verifying is_service_default = true. When false (default), Terraform
  # removes the resource from state only and leaves the configuration in place.
  restore_defaults_on_destroy = true

  depends_on = [time_sleep.wait_30_seconds]

  # --------------------------------------------------------------------------
  # B2B Collaboration Inbound
  # Controls inbound B2B guest access from external tenants into this tenant.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Collaboration Outbound
  # Controls which users in this tenant can be invited as guests to external
  # tenants via B2B collaboration.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Direct Connect Inbound
  # Controls inbound Teams Connect shared channels from external tenants.
  # Note: users_and_groups only supports "AllUsers" — individual user or
  # group GUIDs are not valid for this direction.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Direct Connect Outbound
  # Controls which users in this tenant can use Teams Connect shared channels
  # with external tenants. Supports specific user and group GUIDs as targets.
  # Note: "AllUsers" cannot be mixed with specific user or group GUIDs in the
  # same targets set.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # Inbound Trust
  # Determines which claims from external tenants are honoured when evaluating
  # Conditional Access policies for inbound B2B users.
  # --------------------------------------------------------------------------
  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }

  # --------------------------------------------------------------------------
  # Invitation Redemption Identity Provider Configuration
  # Sets the ordered list of identity providers tried when a B2B guest redeems
  # an invitation. The API normalises this list to the canonical enum order
  # regardless of input, so the order here is informational only.
  # --------------------------------------------------------------------------
  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
  }

  # --------------------------------------------------------------------------
  # Tenant Restrictions
  # Prevents users on managed devices from signing in to other tenants.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # Automatic User Consent Settings
  # Controls whether users can automatically consent to cross-tenant
  # applications without requiring admin approval.
  # --------------------------------------------------------------------------
  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
