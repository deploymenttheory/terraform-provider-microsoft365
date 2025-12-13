# CAU016: Block High User Risk
# Blocks access for high user risk.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau016_block_high_user_risk" {
  display_name = "CAU016-All: Block access for High Risk Users for All Users when Browser and Modern Auth Clients-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]
    user_risk_levels = ["high"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cau016_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau016_exclude.id
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

