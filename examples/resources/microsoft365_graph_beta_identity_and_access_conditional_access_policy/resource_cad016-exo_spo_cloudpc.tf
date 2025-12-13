# CAD016: Token Protection for EXO/SPO/CloudPC on Windows
# Requires token protection for Exchange Online, SharePoint Online, and Cloud PC on Windows.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad016_token_protection_windows" {
  display_name = "CAD016-EXO_SPO_CloudPC: Require token protection when Modern Auth Clients on Windows-v1.2"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

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
        data.microsoft365_graph_beta_applications_service_principal.azure_virtual_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_365.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_cloud_login.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.office_365_exchange_online.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.office_365_sharepoint_online.items[0].app_id,
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

