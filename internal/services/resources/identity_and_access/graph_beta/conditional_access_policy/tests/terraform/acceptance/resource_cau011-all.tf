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

resource "microsoft365_graph_beta_groups_group" "cau011_exclude" {
  display_name     = "EID_UA_CAU011_Exclude"
  mail_nickname    = "eid-ua-cau011-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for CA policy CAU011"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "modern_workplace" {
  display_name     = "EID_UG_ModernWorkplace"
  mail_nickname    = "eid-ug-modernworkplace"
  mail_enabled     = false
  security_enabled = true
  description      = "Members of this group get access to the Modern Workplace"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU011: Block Unlicensed Users
# Blocks access for all users except those who are licensed (e.g., assigned to
# license groups). Useful for enforcing license compliance.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau011_block_unlicensed" {
  display_name = "acc-test-cau011-all: Block access for All users except licensed when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau011_exclude.id,
        microsoft365_graph_beta_groups_group.modern_workplace.id
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

