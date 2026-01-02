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

# Windows Azure Active Directory
data "microsoft365_graph_beta_applications_service_principal" "windows_azure_active_directory" {
  filter_type  = "display_name"
  filter_value = "Windows Azure Active Directory"
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

resource "microsoft365_graph_beta_groups_group" "cau001_exclude" {
  display_name     = "EID_UA_CAU001_EXCLUDE"
  mail_nickname    = "eid-ua-cau001-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU001_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU001A: Require MFA for Guests - Windows Azure AD
# Requires MFA for guest/external users accessing Windows Azure Active Directory.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau001a_guest_mfa_azure_ad" {
  display_name = "acc-test-cau001a-windows_azure_active_directory: Grant Require MFA for guests when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau001_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      include_guests_or_external_users = {
        guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "internalGuest", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = [data.microsoft365_graph_beta_applications_service_principal.windows_azure_active_directory.items[0].app_id] # Windows Azure Active Directory
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

