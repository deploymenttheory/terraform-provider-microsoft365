# CAP004: Block Authentication Transfer
# Blocks authentication transfer methods to prevent token theft.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap004_block_auth_transfer" {
  display_name = "CAP004-All: Block authentication transfer-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cap004_exclude.id
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

    authentication_flows = {
      transfer_methods = "authenticationTransfer"
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

