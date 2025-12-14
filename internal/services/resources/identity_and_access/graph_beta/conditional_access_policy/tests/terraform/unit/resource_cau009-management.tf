# CAU009: Require MFA for Admin Portals
# Requires MFA when accessing admin portals.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau009_admin_portals_mfa" {
  display_name = "CAU009-Management: Grant Require MFA for Admin Portals for All Users when Browser and Modern Auth Clients-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "MicrosoftAdminPortals",
        "66666666-6666-6666-6666-666666666666"
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

