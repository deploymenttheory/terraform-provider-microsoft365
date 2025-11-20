# CAU008: Phishing-Resistant MFA for Admins
# Requires phishing-resistant MFA for admin roles.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau008_admin_phishing_resistant_mfa" {
  display_name = "CAU008-All: Grant Require Phishing Resistant MFA for Admins when Browser and Modern Auth Clients-v1.4"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = [
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555",
        "55555555-5555-5555-5555-555555555555"
      ]
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
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
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000004" # Maps to "phishing_resistant_mfa"
    }
  }
}

