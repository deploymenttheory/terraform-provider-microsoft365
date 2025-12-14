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

resource "microsoft365_graph_beta_groups_group" "cap001_exclude" {
  display_name     = "EID_UA_CAP001_EXCLUDE"
  mail_nickname    = "eid-ua-cap001-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAP001_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAP001: Block Legacy Authentication
# Blocks legacy authentication protocols for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap001_block_legacy_auth" {
  display_name = "acc-test-cap001-all: Block Legacy Authentication for All users when OtherClients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["other"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cap001_exclude.id
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

