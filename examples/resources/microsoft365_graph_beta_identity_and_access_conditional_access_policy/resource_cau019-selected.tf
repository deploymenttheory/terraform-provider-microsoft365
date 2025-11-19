# CAU019: Allow Only Approved Apps for Guests
# Blocks access to all applications for guests except approved apps. This is the
# inverse of CAU003 - allows specific approved apps for guest users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau019_allow_only_approved_apps_guests" {
  display_name = "CAU019-Selected: Only allow approved apps for guests when Browser and Modern Auth Clients-v1.0"
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
        data.microsoft365_graph_beta_applications_service_principal.windows_azure_active_directory.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_approval_management.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.aad_reporting.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.azure_credential_configuration_endpoint_service.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_app_access_panel.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.my_profile.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.my_apps.items[0].app_id,
        // TODOs: both of these id's don't exist in my tenant. probably need to en able a service first
        // and they will appear. 
        //"19db86c3-b2b9-44cc-b339-36da233a3be2", # my sign-ins
        //"4660504c-45b3-4674-a709-71951a6b0763", # Microsoft Invitation Acceptance Portal
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

