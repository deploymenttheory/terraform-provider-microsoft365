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

resource "microsoft365_graph_beta_groups_group" "cau015_exclude" {
  display_name     = "EID_UA_CAU015_EXCLUDE"
  mail_nickname    = "eid-ua-cau015-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU015_EXCLUDE"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cau015_include" {
  display_name     = "EID_UA_CAU015_INCLUDE"
  mail_nickname    = "eid-ua-cau015-include"
  mail_enabled     = false
  security_enabled = true
  description      = "inclusion group for CA policy CAU015_INCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU015: Block High Sign-in Risk
# Blocks access for high sign-in risk.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau015_block_high_signin_risk" {
  display_name = "acc-test-cau015-all: Block access for High Risk Sign-in for All Users when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types    = ["browser", "mobileAppsAndDesktopClients"]
    sign_in_risk_levels = ["high"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cau015_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau015_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

