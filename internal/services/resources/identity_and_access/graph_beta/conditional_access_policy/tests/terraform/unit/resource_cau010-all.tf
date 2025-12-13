# CAU010: Require Terms of Use
# Requires acceptance of terms of use for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau010_terms_of_use" {
  display_name = "CAU010-All: Grant Require ToU for All Users when Browser and Modern Auth Clients-v1.2"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

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
    built_in_controls             = []
    custom_authentication_factors = []
    # include your terms of use ID here
    terms_of_use = ["99999999-9999-9999-9999-999999999999"]
  }
}

