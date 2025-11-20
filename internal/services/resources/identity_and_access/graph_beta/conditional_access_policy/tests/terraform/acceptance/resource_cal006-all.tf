# CAL006: Allow Access Only from Specified Locations
# Restricts access to only specified trusted locations for specific accounts.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal006_allow_only_specified_locations" {
  display_name = "CAL006-All: Only Allow Access from specified locations for specific accounts when Browser and Modern Auth Clients-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cal006_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal006_exclude.id
      ]
      include_roles = []
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
      # Note: Add specific allowed location IDs to exclude_locations
      exclude_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.allowed_apac_office_only.id,
        microsoft365_graph_beta_identity_and_access_named_location.allowed_emea_office_only.id,
        microsoft365_graph_beta_identity_and_access_named_location.allowed_hazelwood_office_only.id,
      ]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

