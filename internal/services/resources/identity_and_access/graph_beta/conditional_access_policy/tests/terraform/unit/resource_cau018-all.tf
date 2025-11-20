# CAU018: Disable Browser Persistence for Admins
# Disables persistent browser sessions for admin users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau018_admin_disable_browser_persistence" {
  display_name = "CAU018-All: Session disable browser persistence for Admins when Browser-v1.0"
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
        "55555555-5555-5555-5555-555555555555"
      ]
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

  session_controls = {
    persistent_browser = {
      mode       = "never"
      is_enabled = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

