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

resource "microsoft365_graph_beta_groups_group" "cad001_exclude" {
  display_name     = "EID_UA_CAD001_Exclude"
  mail_nickname    = "eid-ua-cad001-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for CA policy CAD001"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================

# CAD001: macOS Device Compliance
# Grants macOS access to Office 365 for all users when using modern auth clients
# and device is compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad001_macos_compliant" {
  display_name = "acc-test-cad001-o365: Grant macOS access for All users when Modern Auth Clients and Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad001_exclude.id
      ]
      include_roles = []
      exclude_roles = []
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["macOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice"]
    custom_authentication_factors = []
  }
}

