# CAU017: Admin Sign-in Frequency
# Sets sign-in frequency for admin users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau017_admin_signin_frequency" {
  display_name = "CAU017-All: Session set Sign-in Frequency for Admins when Browser-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

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

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      value               = 10
      type                = "hours"
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

