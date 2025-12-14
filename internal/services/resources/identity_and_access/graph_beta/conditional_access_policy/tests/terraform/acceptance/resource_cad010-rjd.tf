# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

# Break Glass Emergency Access Accounts
resource "microsoft365_graph_beta_groups_group" "breakglass" {
  display_name     = "EID_UA_ConAcc-Breakglass"
  mail_nickname    = "eid-ua-conacc-breakglass"
  mail_enabled     = false
  security_enabled = true
  description      = "Group containing Break Glass Accounts"
}

resource "microsoft365_graph_beta_groups_group" "cad010_exclude" {
  display_name     = "EID_UA_CAD010_EXCLUDE"
  mail_nickname    = "eid-ua-cad010-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAD010_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD010: Require MFA for Device Registration/Join
# Requires MFA when users register or join devices to Entra ID.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad010_device_registration_mfa" {
  display_name = "acc-test-cad010-rjd: Require MFA for device join or registration when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad010_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = []
      exclude_applications                            = []
      include_user_actions                            = ["urn:user:registerdevice"]
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

