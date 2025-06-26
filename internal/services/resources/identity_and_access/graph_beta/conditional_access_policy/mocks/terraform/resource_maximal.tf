resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "maximal" {
  display_name = "Comprehensive Security Policy - Maximal"
  state        = "enabled"

  conditions = {
    client_app_types    = ["all"]
    user_risk_levels    = ["high"]
    sign_in_risk_levels = ["high", "medium"]

    applications = {
      include_applications = ["All"]
      exclude_applications = ["00000002-0000-0ff1-ce00-000000000000"]
      include_user_actions = []
      application_filter = {
        mode = "exclude"
        rule = "device.deviceOwnership -eq \"Company\""
      }
    }

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = []
      exclude_groups = []
      include_roles  = []
      exclude_roles  = ["62e90394-69f5-4237-9190-012177145e10"]
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = ["AllTrusted"]
    }

    devices = {
      include_devices       = []
      exclude_devices       = []
      include_device_states = []
      exclude_device_states = []
      device_filter = {
        mode = "include"
        rule = "device.isCompliant -eq True"
      }
    }
  }

  grant_controls = {
    operator                      = "AND"
    built_in_controls             = ["mfa", "compliantDevice"]
    custom_authentication_factors = []
    terms_of_use                  = []
    authentication_strength = {
      id                     = "00000000-0000-0000-0000-000000000004"
      display_name           = "Multifactor authentication"
      description            = "Combinations of methods that satisfy strong authentication, such as a password + SMS"
      policy_type            = "builtIn"
      requirements_satisfied = "mfa"
      allowed_combinations = [
        "password,sms",
        "password,voice",
        "password,hardwareOath",
        "password,softwareOath",
        "password,microsoftAuthenticatorPush"
      ]
    }
  }

  session_controls = {
    disable_resilience_defaults = false
    application_enforced_restrictions = {
      is_enabled = true
    }
    cloud_app_security = {
      is_enabled              = true
      cloud_app_security_type = "monitorOnly"
    }
    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 4
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
    }
    persistent_browser = {
      is_enabled = false
      mode       = "never"
    }
    continuous_access_evaluation = {
      mode = "strict"
    }
    secure_sign_in_session = {
      is_enabled = true
    }
  }
} 