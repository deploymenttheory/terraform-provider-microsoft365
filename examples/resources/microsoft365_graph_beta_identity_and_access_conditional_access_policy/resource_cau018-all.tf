# CAU018: Disable Browser Persistence for Admins
# Disables persistent browser sessions for admin users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau018_admin_disable_browser_persistence" {
  display_name = "CAU018-All: Session disable browser persistence for Admins when Browser-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau018_exclude.id
      ]
      include_roles = [
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_administrator.items[0].id,
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
    persistent_browser = {
      mode       = "never"
      is_enabled = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

