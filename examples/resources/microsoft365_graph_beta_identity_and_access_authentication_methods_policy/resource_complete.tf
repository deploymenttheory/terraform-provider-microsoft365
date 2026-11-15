# Manage the policy-level settings and every supported method configuration.
resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "this" {
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
      exclude_targets = []
    }
  }

  report_suspicious_activity_settings = {
    state                = "enabled"
    voice_reporting_code = 0
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

  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  temporary_access_pass = {
    state                       = "enabled"
    default_lifetime_in_minutes = 60
    default_length              = 8
    minimum_lifetime_in_minutes = 60
    maximum_lifetime_in_minutes = 480
    is_usable_once              = false
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  sms = {
    state = "disabled"
    include_targets = [
      {
        id                    = "all_users"
        target_type           = "group"
        is_usable_for_sign_in = false
      }
    ]
    exclude_targets = []
  }

  software_oath = {
    state = "enabled"
    include_targets = [
      {
        id          = "all_users"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }

  hardware_oath = {
    state = "disabled"
  }
}
