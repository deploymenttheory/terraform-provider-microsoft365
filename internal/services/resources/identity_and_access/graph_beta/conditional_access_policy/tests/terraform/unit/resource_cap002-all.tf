# CAP002: Block Exchange ActiveSync
# Blocks Exchange ActiveSync clients for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap002_block_exchange_activesync" {
  display_name = "CAP002-All: Block Exchange ActiveSync Clients for All users-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["exchangeActiveSync"]

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

