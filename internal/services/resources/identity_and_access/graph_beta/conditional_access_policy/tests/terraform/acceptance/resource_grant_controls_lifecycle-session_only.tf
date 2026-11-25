# Grant controls lifecycle: a session-only policy (no grant_controls).
# Graph stores a grant with no control as null, so a session-only policy leaves grant_controls out.

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "grant_controls_lifecycle" {
  display_name = "acc-test-grant-controls-lifecycle: Session-only policy for no users ${random_string.suffix.result}"
  state        = "disabled"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["None"]
      exclude_users  = []
      include_groups = []
      exclude_groups = []
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    application_enforced_restrictions = {
      is_enabled = true
    }
  }
}
