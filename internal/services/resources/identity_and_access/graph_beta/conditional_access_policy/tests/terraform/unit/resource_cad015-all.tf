# CAD015: Compliant Device for Windows and macOS Browser Access
# Requires compliant device for all users accessing all apps via browser on Windows/macOS.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad015_windows_macos_browser_compliant" {
  display_name = "CAD015-All: Grant access for All users when Browser and Modern Auth Clients and Compliant on Windows and macOS-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = ["77777777-7777-7777-7777-777777777777"]
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
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

    platforms = {
      include_platforms = ["windows", "macOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }
}

