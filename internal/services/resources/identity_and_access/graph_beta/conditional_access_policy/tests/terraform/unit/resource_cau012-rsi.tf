# CAU012: Security Info Registration with TAP
# Requires MFA for combined security info registration and sets sign-in frequency
# to every time when registering from non-trusted locations. Supports Temporary Access Pass (TAP).
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau012_security_info_registration_tap" {
  display_name = "CAU012-RSI: Combined Security Info Registration with TAP-v1.1"
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

