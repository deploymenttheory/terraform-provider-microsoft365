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

resource "microsoft365_graph_beta_groups_group" "cad016_exclude" {
  display_name     = "EID_UA_CAD016_EXCLUDE"
  mail_nickname    = "eid-ua-cad016-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAD016_EXCLUDE"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cad016_include" {
  display_name     = "EID_UA_CAD016_INCLUDE"
  mail_nickname    = "eid-ua-cad016-include"
  mail_enabled     = false
  security_enabled = true
  description      = "inclusion group for CA policy CAD016_INCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD016: Token Protection for EXO/SPO/CloudPC on Windows
# Requires token protection for Exchange Online, SharePoint Online, and Cloud PC on Windows.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad016_token_protection_windows" {
  display_name = "acc-test-cad016-exo_spo_cloudpc: Require token protection when Modern Auth Clients on Windows ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad016_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad016_exclude.id
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
      include_applications = [
        "9cdead84-a844-4324-93f2-b2e6bb768d07",
        "0af06dc6-e4b5-4f28-818e-e78e62d137a5",
        "270efc09-cd0d-444b-a71f-39af4910ec45",
        "00000002-0000-0ff1-ce00-000000000000",
        "00000003-0000-0ff1-ce00-000000000000",
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

