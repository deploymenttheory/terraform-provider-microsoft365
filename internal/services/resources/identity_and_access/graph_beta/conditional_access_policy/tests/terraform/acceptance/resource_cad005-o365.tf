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

resource "microsoft365_graph_beta_groups_group" "cad005_exclude" {
  display_name     = "EID_UA_CAD005_EXCLUDE"
  mail_nickname    = "eid-ua-cad005-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAD005_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD005: Block Unsupported Device Platforms
# Blocks access to Office 365 for unsupported device platforms.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad005_block_unsupported_platforms" {
  display_name = "acc-test-cad005-o365: Block access for unsupported device platforms for All users when Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad005_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = ["android", "iOS", "windows", "macOS", "linux"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

