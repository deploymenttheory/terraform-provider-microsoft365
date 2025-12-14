# CAU020: Insider Risk Conditional Access Policy
# Block access for Elevated Insider Risk Users for all Users
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau020_insider_risk_block" {
  display_name = "CAU020-ALL: Block access for Elevated Insider Risk Users for all Users v1.0"
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
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels           = []
    user_risk_levels              = []
    service_principal_risk_levels = []
    insider_risk_levels           = ["moderate", "elevated"]
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}
