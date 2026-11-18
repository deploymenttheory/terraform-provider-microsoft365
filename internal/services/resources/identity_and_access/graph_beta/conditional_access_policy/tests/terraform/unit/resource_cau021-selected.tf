# CAU021: Sign-in Every Time Strict for Private Access Apps
# Requires reauthentication on each new connection to selected Global Secure Access
# Private Access apps. Uses the "Every time strict (preview)" sign-in frequency.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau021_private_access_sign_in_every_time" {
  display_name = "CAU021-Selected: Session require sign-in every time strict for Private Access apps-v1.0"
  state        = "disabled"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = ["44444444-4444-4444-4444-444444444444"]
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["55555555-5555-5555-5555-555555555555"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTimeStrict"
      is_enabled          = true
    }
  }
}
