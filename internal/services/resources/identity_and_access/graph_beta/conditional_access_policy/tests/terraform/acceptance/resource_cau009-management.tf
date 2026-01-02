# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Data Sources
# ==============================================================================

# Windows Azure Service Management API
data "microsoft365_graph_beta_applications_service_principal" "windows_azure_service_management_api" {
  filter_type  = "display_name"
  filter_value = "Windows Azure Service Management API"
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

resource "microsoft365_graph_beta_groups_group" "cau009_exclude" {
  display_name     = "EID_UA_CAU009_EXCLUDE"
  mail_nickname    = "eid-ua-cau009-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU009_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU009: Require MFA for Admin Portals
# Requires MFA when accessing admin portals.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau009_admin_portals_mfa" {
  display_name = "acc-test-cau009-management: Grant Require MFA for Admin Portals for All Users when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau009_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "MicrosoftAdminPortals",
        data.microsoft365_graph_beta_applications_service_principal.windows_azure_service_management_api.items[0].app_id
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

