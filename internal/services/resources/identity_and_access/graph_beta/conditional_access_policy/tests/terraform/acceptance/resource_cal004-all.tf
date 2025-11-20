# CAL004: Block Admin Access from Non-Trusted Locations
# Blocks admin access from non-trusted locations.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal004_block_admin_untrusted_locations" {
  display_name = "CAL004-All: Block access for Admins from non-trusted locations when Browser and Modern Auth Clients-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal004_exclude.id
      ]
      include_roles = [
        data.microsoft365_graph_beta_identity_and_access_role_definitions.global_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.billing_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_operator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.conditional_access_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_device_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.user_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.compliance_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.directory_writers.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.exchange_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.global_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.helpdesk_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.hybrid_identity_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.insights_business_leader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.intune_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.knowledge_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_role_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.reports_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.search_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.sharepoint_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.teams_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_administrator.items[0].id,
      ]
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = ["AllTrusted"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

