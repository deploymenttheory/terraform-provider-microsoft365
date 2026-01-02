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

resource "microsoft365_graph_beta_groups_group" "cau012_exclude" {
  display_name     = "EID_UA_CAU012_EXCLUDE"
  mail_nickname    = "eid-ua-cau012-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU012_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU012: Security Info Registration with TAP
# Requires MFA for combined security info registration and sets sign-in frequency
# to every time when registering from non-trusted locations. Supports Temporary Access Pass (TAP).
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau012_security_info_registration_tap" {
  display_name = "acc-test-cau012-rsi: Combined Security Info Registration with TAP ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau012_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = []
      exclude_applications                            = []
      include_user_actions                            = ["urn:user:registersecurityinfo"]
      include_authentication_context_class_references = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = ["AllTrusted"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
    }
  }
}

