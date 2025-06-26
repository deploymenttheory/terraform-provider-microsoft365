resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "minimal" {
  display_name = "Block Legacy Authentication - Minimal"
  state        = "enabled"

  conditions = {
    client_app_types    = ["exchangeActiveSync", "other"]
    user_risk_levels    = []
    sign_in_risk_levels = []

    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
    }

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = []
    }
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["block"]
  }
} 