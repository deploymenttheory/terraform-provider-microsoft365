# CAU011: Block Unlicensed Users
# Blocks access for all users except those who are licensed (e.g., assigned to
# license groups). Useful for enforcing license compliance.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau011_block_unlicensed" {
  display_name = "CAU011-All: Block access for All users except licensed when Browser and Modern Auth Clients-v1.0"
  state        = "disabled" # Note: Original policy was disabled

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau011_exclude.id,
        microsoft365_graph_beta_groups_group.modern_workplace.id
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

