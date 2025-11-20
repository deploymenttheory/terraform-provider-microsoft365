# CAP003: Block Device Code Flow
# Blocks device code authentication flow to prevent phishing attacks.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap003_block_device_code_flow" {
  display_name = "CAP003-All: Block device code authentication flow-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

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

    authentication_flows = {
      transfer_methods = "deviceCodeFlow"
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

