# CAL001: Block Specified Locations
# Blocks access from specified untrusted locations for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal001_block_locations" {
  display_name = "CAL001-All: Block specified locations for All users when Browser and Modern Auth Clients-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal001_exclude.id
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
      # Note: Add specific blocked location IDs
      include_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_client_ip.id,
        microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_authenticator_gps.id
      ]
      exclude_locations = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

