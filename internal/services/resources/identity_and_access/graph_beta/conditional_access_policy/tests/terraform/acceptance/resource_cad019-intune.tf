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
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cad019_exclude" {
  display_name     = "EID_UA_CAD019_EXCLUDE"
  mail_nickname    = "eid-ua-cad019-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAD019_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD019: Intune Enrollment - MFA and Sign-in Frequency
# Requires MFA and sets sign-in frequency to every time for Intune enrollment.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad019_intune_enrollment_mfa" {
  display_name = "acc-test-cad019-intune: Require MFA and set sign-in frequency to every time ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad019_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "d4ebce55-015a-49b5-a083-c84d1797ae8c" # Microsoft Intune Enrollment
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # multifactor_authentication
    }
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
      # Note: type and value are not set when frequency_interval is "everyTime"
    }
  }
}

