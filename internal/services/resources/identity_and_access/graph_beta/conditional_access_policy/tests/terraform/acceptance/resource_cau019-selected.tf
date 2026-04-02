# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Application Dependencies
# These are well-known Microsoft first-party applications. Using data source
# lookups ensures the appId (client ID) is retrieved from the tenant rather
# than relying on hardcoded values that may not be provisioned.
# ==============================================================================

data "microsoft365_graph_beta_applications_service_principal" "windows_azure_active_directory" {
  display_name = "Windows Azure Active Directory"
}

data "microsoft365_graph_beta_applications_service_principal" "microsoft_approval_management" {
  display_name = "Microsoft Approval Management"
}

data "microsoft365_graph_beta_applications_service_principal" "aad_reporting" {
  display_name = "Azure Active Directory Graph"
}

data "microsoft365_graph_beta_applications_service_principal" "azure_credential_configuration_endpoint_service" {
  display_name = "Azure Credential Configuration Endpoint Service"
}

data "microsoft365_graph_beta_applications_service_principal" "microsoft_app_access_panel" {
  display_name = "Microsoft App Access Panel"
}

data "microsoft365_graph_beta_applications_service_principal" "my_profile" {
  display_name = "My Profile"
}

data "microsoft365_graph_beta_applications_service_principal" "my_apps" {
  display_name = "My Apps"
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

resource "microsoft365_graph_beta_groups_group" "cau019_exclude" {
  display_name     = "EID_UA_CAU019_EXCLUDE"
  mail_nickname    = "eid-ua-cau019-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAU019_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU019: Allow Only Approved Apps for Guests
# Blocks access to all applications for guests except approved apps. This is the
# inverse of CAU003 - allows specific approved apps for guest users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau019_allow_only_approved_apps_guests" {
  display_name = "acc-test-cau019-selected: Only allow approved apps for guests when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau019_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      include_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications = ["All"]
      exclude_applications = [
        data.microsoft365_graph_beta_applications_service_principal.windows_azure_active_directory.app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_approval_management.app_id,
        data.microsoft365_graph_beta_applications_service_principal.aad_reporting.app_id,
        data.microsoft365_graph_beta_applications_service_principal.azure_credential_configuration_endpoint_service.app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_app_access_panel.app_id,
        data.microsoft365_graph_beta_applications_service_principal.my_profile.app_id,
        data.microsoft365_graph_beta_applications_service_principal.my_apps.app_id,
        "Office365"
      ]
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

