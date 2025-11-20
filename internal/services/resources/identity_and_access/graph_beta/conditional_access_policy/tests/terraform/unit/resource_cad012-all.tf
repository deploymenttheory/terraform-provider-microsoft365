# CAD012: Admin Access on Compliant Devices
# Grants access for admin users to all apps when using compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad012_admin_compliant_access" {
  display_name = "CAD012-All: Grant access for Admin users when Browser and Modern Auth Clients and Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = [
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
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

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }
}

