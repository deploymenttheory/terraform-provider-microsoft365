# CAU017: Admin Sign-in Frequency
# Sets sign-in frequency for admin users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau017_admin_signin_frequency" {
  display_name = "CAU017-All: Session set Sign-in Frequency for Admins when Browser-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau017_exclude.id
      ]
      include_roles = [
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_developer.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.authentication_extensibility_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.b2c_ief_keyset_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.billing_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_device_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.compliance_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.conditional_access_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.directory_writers.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.exchange_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.global_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.global_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.helpdesk_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.hybrid_identity_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.intune_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.password_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_role_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_operator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.sharepoint_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.teams_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.user_administrator.items[0].id
      ]
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      value               = 10
      type                = "hours"
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

