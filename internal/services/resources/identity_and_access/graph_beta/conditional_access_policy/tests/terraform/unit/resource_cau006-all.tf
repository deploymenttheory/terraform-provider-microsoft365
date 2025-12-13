# CAU006: MFA for Medium/High Sign-in Risk
# Requires MFA for medium and high sign-in risk.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau006_signin_risk_mfa" {
  display_name = "CAU006-All: Grant access for Medium and High Risk Sign-in for All Users when Browser and Modern Auth Clients require MFA-v1.4"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types    = ["browser", "mobileAppsAndDesktopClients"]
    sign_in_risk_levels = ["high", "medium"]

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
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
    }
  }
}

