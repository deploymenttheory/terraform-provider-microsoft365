resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "minimal" {
  display_name = "Block Legacy Authentication - Minimal"
  state        = "enabled"

  conditions = {
    client_app_types    = ["exchangeActiveSync", "other"]
    sign_in_risk_levels = []  # Required field but can be empty

    applications = {
      include_applications = ["All"]
    }

    users = {
      include_users = ["All"]
    }
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}