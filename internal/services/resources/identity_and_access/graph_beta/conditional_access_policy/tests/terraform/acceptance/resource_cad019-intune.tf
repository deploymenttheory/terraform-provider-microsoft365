# CAD019: Intune Enrollment - MFA and Sign-in Frequency
# Requires MFA and sets sign-in frequency to every time for Intune enrollment.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad019_intune_enrollment_mfa" {
  display_name = "CAD019-Intune: Require MFA and set sign-in frequency to every time-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad019_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        data.microsoft365_graph_beta_applications_service_principal.microsoft_intune_enrollment.items[0].app_id
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # multifactor_authentication
    }
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
      # Note: type and value are not set when frequency_interval is "everyTime"
    }
  }
}

