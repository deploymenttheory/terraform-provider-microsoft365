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

resource "microsoft365_graph_beta_groups_group" "cau003_exclude" {
  display_name     = "EID_UA_CAU003_EXCLUDE"
  mail_nickname    = "eid-ua-cau003-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU003_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Application Dependencies - chosen because it's built-in and will always be available
# ==============================================================================

# Windows Azure Service Management API - built-in application
data "microsoft365_graph_beta_applications_service_principal" "windows_azure_service_management_api" {
  filter_type  = "display_name"
  filter_value = "Windows Azure Service Management API"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU003: Block Unapproved Apps for Guests
# Blocks access to unapproved applications for guest users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau003_block_unapproved_apps_guests" {
  display_name = "acc-test-cau003-selected: Block unapproved apps for guests when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau003_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      include_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = [data.microsoft365_graph_beta_applications_service_principal.windows_azure_service_management_api.items[0].app_id]
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

