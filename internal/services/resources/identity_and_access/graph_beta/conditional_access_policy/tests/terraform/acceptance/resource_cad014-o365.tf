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

resource "microsoft365_graph_beta_groups_group" "cad014_exclude" {
  display_name     = "EID_UA_CAD014_EXCLUDE"
  mail_nickname    = "eid-ua-cad014-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAD014_EXCLUDE"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cad014_include" {
  display_name     = "EID_UA_CAD014_INCLUDE"
  mail_nickname    = "eid-ua-cad014-include"
  mail_enabled     = false
  security_enabled = true
  description      = "uincludeion group for CA policy CAD014_INCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD014: Edge App Protection on Windows
# Requires app protection policy for Edge browser on Windows for Office 365 access.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad014_edge_app_protection_windows" {
  display_name = "acc-test-cad014-o365: Require App Protection Policy for Edge on Windows for All users when Browser and Non-Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad014_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad014_exclude.id
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
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantApplication"]
    custom_authentication_factors = []
  }
}

