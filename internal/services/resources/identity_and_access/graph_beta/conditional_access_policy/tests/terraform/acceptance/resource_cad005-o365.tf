# CAD005: Block Unsupported Device Platforms
# Blocks access to Office 365 for unsupported device platforms.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad005_block_unsupported_platforms" {
  display_name = "CAD005-O365: Block access for unsupported device platforms for All users when Modern Auth Clients-v1.1"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad005_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = ["android", "iOS", "windows", "macOS", "linux"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

