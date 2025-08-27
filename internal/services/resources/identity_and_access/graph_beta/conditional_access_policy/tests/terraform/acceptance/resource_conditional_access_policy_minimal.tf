resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "minimal" {
  display_name = "tf-test"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    applications = {
      include_applications                             = ["All"]
      exclude_applications                             = []
      include_user_actions                             = []
      include_authentication_context_class_references = []
    }

    users = {
      include_users  = ["None"]
      exclude_users  = []
      include_groups = []
      exclude_groups = []
      include_roles  = []
      exclude_roles  = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = []
    }

    sign_in_risk_levels = []
    user_risk_levels    = []
  }

  grant_controls = {
    operator                        = "OR"
    built_in_controls              = []
    custom_authentication_factors = []
    terms_of_use                   = ["79f28780-c502-49c4-8951-f53f6a239b60"]
  }
}