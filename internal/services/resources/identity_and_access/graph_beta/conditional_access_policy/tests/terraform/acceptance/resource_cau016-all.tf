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

resource "microsoft365_graph_beta_groups_group" "cau016_exclude" {
  display_name     = "EID_UA_CAU016_EXCLUDE"
  mail_nickname    = "eid-ua-cau016-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU016_EXCLUDE"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cau016_include" {
  display_name     = "EID_UA_CAU016_INCLUDE"
  mail_nickname    = "eid-ua-cau016-include"
  mail_enabled     = false
  security_enabled = true
  description      = "inclusion group for CA policy CAU016_INCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU016: Block High User Risk
# Blocks access for high user risk.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau016_block_high_user_risk" {
  display_name = "acc-test-cau016-all: Block access for High Risk Users for All Users when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]
    user_risk_levels = ["high"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cau016_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau016_exclude.id
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

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

