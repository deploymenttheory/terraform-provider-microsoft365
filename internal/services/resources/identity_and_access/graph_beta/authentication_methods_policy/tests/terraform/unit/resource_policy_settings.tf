resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  policy_migration_state = "migrationComplete"

  registration_enforcement = {
    authentication_methods_registration_campaign = {
      snooze_duration_in_days                    = 3
      enforce_registration_after_allowed_snoozes = true
      state                                      = "enabled"
      include_targets = [
        {
          id                             = "all_users"
          target_type                    = "group"
          targeted_authentication_method = "microsoftAuthenticator"
        }
      ]
      exclude_targets = [
        {
          id          = "33333333-3333-3333-3333-333333333333"
          target_type = "group"
        }
      ]
    }
  }

  report_suspicious_activity_settings = {
    state                = "enabled"
    voice_reporting_code = 8
    include_target = {
      id          = "all_users"
      target_type = "group"
    }
  }

  system_credential_preferences = {
    state = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }
}
