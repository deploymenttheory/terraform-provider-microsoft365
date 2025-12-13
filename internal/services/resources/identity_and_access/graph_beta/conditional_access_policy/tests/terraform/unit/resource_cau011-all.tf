# CAU011: Block Unlicensed Users
# Blocks access for all users except those who are licensed (e.g., assigned to
# license groups). Useful for enforcing license compliance.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau011_block_unlicensed" {
  display_name = "CAU011-All: Block access for All users except licensed when Browser and Modern Auth Clients-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333",
        "88888888-8888-8888-8888-888888888888"
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

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

