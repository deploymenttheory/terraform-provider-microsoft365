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

resource "microsoft365_graph_beta_groups_group" "cau020_exclude" {
  display_name     = "EID_UA_CAU020_Exclude"
  mail_nickname    = "eid-ua-cau020-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for CA policy CAU020"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================

# CAU020: Insider Risk Conditional Access Policy
# Block access for Elevated Insider Risk Users for all Users
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau020_all" {
  display_name = "acc-test-cau020-all: Block access for Elevated Insider Risk Users for all Users ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau020_exclude.id
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

    sign_in_risk_levels           = []
    user_risk_levels              = []
    service_principal_risk_levels = []
    agent_id_risk_levels          = []
    insider_risk_levels           = ["moderate", "elevated"]
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}
