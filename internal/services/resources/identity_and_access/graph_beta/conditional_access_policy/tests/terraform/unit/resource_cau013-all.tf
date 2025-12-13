# CAU013: Phishing-Resistant MFA for All Users
# Requires phishing-resistant MFA for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau013_all_users_phishing_resistant_mfa" {
  display_name = "CAU013-All: Grant Require phishing resistant MFA for All users when Browser and Modern Auth Clients-v1.0"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = ["77777777-7777-7777-7777-777777777777"]
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = ["All"]
      exclude_applications = [
        "66666666-6666-6666-6666-666666666666"
      ]
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000004" # Maps to "phishing_resistant_mfa"
    }
  }
}

