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

resource "microsoft365_graph_beta_groups_group" "cad017_exclude" {
  display_name     = "EID_UA_CAD017_EXCLUDE"
  mail_nickname    = "eid-ua-cad017-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAD017_EXCLUDE"
}

resource "microsoft365_graph_beta_groups_group" "cad017_include" {
  display_name     = "EID_UA_CAD017_INCLUDE"
  mail_nickname    = "eid-ua-cad017-include"
  mail_enabled     = false
  security_enabled = true
  description      = "uincludeion group for CA policy CAD017_INCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD017: Selected Apps - Mobile App Protection or Compliance
# Requires app protection policy or device compliance for selected apps on iOS/Android.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad017_selected_mobile_app_protection" {
  display_name = "acc-test-cad017-selected: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad017_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad017_exclude.id
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
      include_applications                            = ["None"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "compliantApplication"]
    custom_authentication_factors = []
  }
}

