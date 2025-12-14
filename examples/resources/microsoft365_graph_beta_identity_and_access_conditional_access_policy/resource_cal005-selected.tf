# CAL005: Less-Trusted Locations Require Compliance
# Requires compliant device when accessing from less-trusted locations.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal005_less_trusted_locations_compliant" {
  display_name = "CAL005-Selected: Grant access for All users on less-trusted locations when Browser and Modern Auth Clients and Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal005_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = ["Office365"]
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    locations = {
      # Note: Add specific less-trusted location IDs
      include_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_partner_networks.id,
        microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_public_spaces.id
      ]
      exclude_locations = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }
}

