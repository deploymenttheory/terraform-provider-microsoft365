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

resource "microsoft365_graph_beta_groups_group" "cap002_exclude" {
  display_name     = "EID_UA_CAP002_EXCLUDE"
  mail_nickname    = "eid-ua-cap002-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAP002_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAP002: Block Exchange ActiveSync
# Blocks Exchange ActiveSync clients for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap002_block_exchange_activesync" {
  display_name = "acc-test-cap002-all: Block Exchange ActiveSync Clients for All users ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["exchangeActiveSync"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cap002_exclude.id
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

