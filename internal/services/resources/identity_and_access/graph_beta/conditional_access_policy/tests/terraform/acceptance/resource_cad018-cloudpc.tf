# ==============================================================================
# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

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

resource "microsoft365_graph_beta_groups_group" "cad018_exclude" {
  display_name     = "EID_UA_CAD018_EXCLUDE"
  mail_nickname    = "eid-ua-cad018-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAD018_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD018: Cloud PC - Mobile App Protection or Compliance
# Requires app protection policy or device compliance for Cloud PC access on iOS/Android.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad018_cloudpc_mobile_app_protection" {
  display_name = "acc-test-cad018-cloudpc: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad018_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "9cdead84-a844-4324-93f2-b2e6bb768d07",
        "a4a365df-50f1-4397-bc59-1a1564b8bb9c",
        "0af06dc6-e4b5-4f28-818e-e78e62d137a5",
        "270efc09-cd0d-444b-a71f-39af4910ec45"
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "compliantApplication"]
    custom_authentication_factors = []
  }
}

