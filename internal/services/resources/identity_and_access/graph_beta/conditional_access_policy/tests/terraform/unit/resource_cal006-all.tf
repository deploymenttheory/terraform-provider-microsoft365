# CAL006: Allow Access Only from Specified Locations
# Restricts access to only specified trusted locations for specific accounts.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal006_allow_only_specified_locations" {
  display_name = "CAL006-All: Only Allow Access from specified locations for specific accounts when Browser and Modern Auth Clients-v1.0"
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
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    locations = {
      include_locations = ["All"]
      # Note: Add specific allowed location IDs to exclude_locations
      exclude_locations = [
        "44444444-4444-4444-4444-444444444444",
        "44444444-4444-4444-4444-444444444444",
        "44444444-4444-4444-4444-444444444444",
      ]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

