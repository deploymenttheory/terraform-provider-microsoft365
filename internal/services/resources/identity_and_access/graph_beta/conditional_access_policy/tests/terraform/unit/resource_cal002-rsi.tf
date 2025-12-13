# CAL002: MFA Registration from Trusted Locations Only
# Requires security info registration to occur from trusted locations only.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal002_mfa_registration_trusted_locations" {
  display_name = "CAL002-RSI: Require MFA registration from trusted locations only for All users when Browser and Modern Auth Clients-v1.4"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

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
      include_applications                            = []
      exclude_applications                            = []
      include_user_actions                            = ["urn:user:registersecurityinfo"]
      include_authentication_context_class_references = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = ["AllTrusted"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

