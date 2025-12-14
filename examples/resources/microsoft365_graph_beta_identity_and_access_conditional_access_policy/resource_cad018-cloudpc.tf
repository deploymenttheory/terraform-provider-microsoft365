# CAD018: Cloud PC - Mobile App Protection or Compliance
# Requires app protection policy or device compliance for Cloud PC access on iOS/Android.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad018_cloudpc_mobile_app_protection" {
  display_name = "CAD018-CloudPC: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad018_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        data.microsoft365_graph_beta_applications_service_principal.azure_virtual_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_remote_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_365.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_cloud_login.items[0].app_id
      ]
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

