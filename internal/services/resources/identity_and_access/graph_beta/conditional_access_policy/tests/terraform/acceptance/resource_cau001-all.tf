# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Service Principal Dependencies
# ==============================================================================

# Microsoft Rights Management Services
data "microsoft365_graph_beta_applications_service_principal" "microsoft_rights_management_services" {
  filter_type  = "display_name"
  filter_value = "Microsoft Rights Management Services"
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

resource "microsoft365_graph_beta_groups_group" "cau001_exclude" {
  display_name     = "EID_UA_CAU001_EXCLUDE"
  mail_nickname    = "eid-ua-cau001-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU001_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU001: Require MFA for Guest Users
# Requires MFA for guest/external users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau001_guest_mfa" {
  display_name = "acc-test-cau001-all: Grant Require MFA for guests when Browser and Modern Auth Clients ${random_string.suffix.result}"
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
      include_applications                            = ["All"]
      exclude_applications                            = [data.microsoft365_graph_beta_applications_service_principal.microsoft_rights_management_services.items[0].app_id]
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

